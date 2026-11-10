package parser

import (
	"fmt"
	"plain/pkg/ast"
	"plain/pkg/token"
	"strconv"
)

// Precedence rankings for Pratt parsing
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOTEQ:    EQUALS,
	token.LT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GT:       LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	// Register prefix parse functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.NONE, p.parseNoneLiteral)

	// Register infix parse functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOTEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignmentStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional As TypeName clause
	if p.peekTokenIs(token.AS) {
		p.nextToken() // move to AS
		p.nextToken() // move to TypeName
		stmt.TypeName = p.curToken.Literal
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	// Parse the condition expression (e.g., total == 30)
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	// Parse the consequence body until Else or End
	stmt.Consequence = p.parseBlockStatement()

	// Check for an optional Else block
	if p.curToken.Type == token.ELSE {
		stmt.Alternative = p.parseBlockStatement()
	}

	// Must terminate with End
	if p.curToken.Type != token.END {
		p.peekError(token.END)
		return nil
	}

	// Followed by If (End If)
	if !p.expectPeek(token.IF) {
		return nil
	}

	return stmt
}

func (p *Parser) parseBlockStatement() []ast.Statement {
	var statements []ast.Statement
	p.nextToken() // Move past Then or Else

	for p.curToken.Type != token.EOF && p.curToken.Type != token.ELSE && p.curToken.Type != token.END {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.VarName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.AS) {
		return nil
	}

	if !p.expectPeek(token.NUMERIC) {
		return nil
	}
	stmt.TypeName = p.curToken.Literal

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Start = p.parseExpression(LOWEST)

	if !p.expectPeek(token.TO) {
		return nil
	}

	p.nextToken()
	stmt.End = p.parseExpression(LOWEST)

	// Check for optional Step
	if p.peekTokenIs(token.STEP) {
		p.nextToken() // move to STEP
		p.nextToken() // move past STEP
		stmt.Step = p.parseExpression(LOWEST)
	}

	// Parse loop body statements until End
	stmt.Body = p.parseForBlockStatement()

	if p.curToken.Type != token.END {
		p.peekError(token.END)
		return nil
	}

	if !p.expectPeek(token.FOR) {
		return nil
	}

	return stmt
}

func (p *Parser) parseForBlockStatement() []ast.Statement {
	var statements []ast.Statement
	p.nextToken() // Move past To/Step expression end or newline

	for p.curToken.Type != token.EOF && p.curToken.Type != token.END {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	// Optional Return Type: As TypeName
	if p.peekTokenIs(token.AS) {
		p.nextToken() // move to AS
		if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.NUMERIC) || p.peekTokenIs(token.STRING) {
			p.nextToken()
			stmt.ReturnType = p.curToken.Literal
		} else {
			p.peekError(token.IDENT)
			return nil
		}
	}
	stmt.Body = p.parseFunctionBlockStatement()

	if p.curToken.Type != token.END {
		p.peekError(token.END)
		return nil
	}

	if !p.expectPeek(token.FUNCTION) {
		return nil
	}

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	var parameters []*ast.Parameter

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()
	typeName := p.curToken.Literal

	param := &ast.Parameter{Token: ident.Token, Name: ident, TypeName: typeName}
	parameters = append(parameters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if !p.expectPeek(token.AS) {
			return nil
		}
		p.nextToken()
		typeName := p.curToken.Literal
		param := &ast.Parameter{Token: ident.Token, Name: ident, TypeName: typeName}
		parameters = append(parameters, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseFunctionBlockStatement() []ast.Statement {
	var statements []ast.Statement
	p.nextToken() // Move past header / closing paren / return type

	for p.curToken.Type != token.EOF && p.curToken.Type != token.END {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	if p.peekTokenIs(token.END) || p.peekTokenIs(token.EOF) {
		return stmt // Return with no expression (implicitly returns None)
	}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

// Prefix parsers
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as number", p.curToken.Literal))
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// e.g., -10
	return nil // Expand later if needed for PoC
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // consume '('

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// Infix parsers
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curTokenPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Helper methods
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curTokenPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // move to comma
		p.nextToken() // move past comma to next expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
}

func (p *Parser) parseNoneLiteral() ast.Expression {
	return &ast.NoneLiteral{Token: p.curToken}
}
