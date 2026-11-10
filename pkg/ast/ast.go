package ast

import (
	"plain/pkg/token"
	"strings"
)

// Node represents any node in the Abstract Syntax Tree.
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

// --- Statements ---

type VarStatement struct {
	token.Token // The 'Var' token
	Name        *Identifier
	TypeName    string
	Value       Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Literal }
func (vs *VarStatement) String() string {
	return "Var " + vs.Name.String() + " As " + vs.TypeName + " = " + vs.Value.String()
}

type ExpressionStatement struct {
	token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// --- Expressions ---

type Identifier struct {
	token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	token.Token
	Value float64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Literal }
func (il *IntegerLiteral) String() string       { return il.Literal }

type InfixExpression struct {
	token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// IfStatement represents an if-else conditional block
type IfStatement struct {
	token.Token // The 'If' token
	Condition   Expression
	Consequence []Statement
	Alternative []Statement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Literal }
func (is *IfStatement) String() string {
	var out string
	out += "If " + is.Condition.String() + " Then "
	for _, s := range is.Consequence {
		out += s.String()
	}
	if len(is.Alternative) > 0 {
		out += " Else "
		for _, s := range is.Alternative {
			out += s.String()
		}
	}
	out += " End If"
	return out
}

// ForStatement represents a range or iterator loop block
type ForStatement struct {
	token.Token // The 'For' token
	VarName     *Identifier
	TypeName    string
	Start       Expression
	End         Expression
	Step        Expression // Optional; can be nil (defaults to 1)
	Body        []Statement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Literal }
func (fs *ForStatement) String() string {
	var out string
	out += "For " + fs.VarName.String() + " As " + fs.TypeName + " = " + fs.Start.String() + " To " + fs.End.String()
	if fs.Step != nil {
		out += " Step " + fs.Step.String()
	}
	for _, s := range fs.Body {
		out += s.String()
	}
	out += " End For"
	return out
}

type AssignmentStatement struct {
	token.Token // The identifier token
	Name        *Identifier
	Value       Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Literal }
func (as *AssignmentStatement) String() string {
	return as.Name.String() + " = " + as.Value.String()
}

// Parameter represents a function parameter (e.g., a As Numeric)
type Parameter struct {
	token.Token
	Name     *Identifier
	TypeName string
}

func (p *Parameter) statementNode()       {}
func (p *Parameter) TokenLiteral() string { return p.Literal }
func (p *Parameter) String() string {
	return p.Name.String() + " As " + p.TypeName
}

// FunctionStatement represents a function declaration block
type FunctionStatement struct {
	token.Token // The 'Function' token
	Name        *Identifier
	Parameters  []*Parameter
	ReturnType  string // Optional; empty if none
	Body        []Statement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Literal }
func (fs *FunctionStatement) String() string {
	var out string
	out += "Function " + fs.Name.String() + "("
	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}
	out += strings.Join(params, ", ") + ")"
	if fs.ReturnType != "" {
		out += " As " + fs.ReturnType
	}
	for _, s := range fs.Body {
		out += s.String()
	}
	out += " End Function"
	return out
}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	token.Token // The 'Return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Literal }
func (rs *ReturnStatement) String() string {
	var out string
	out += "Return"
	if rs.ReturnValue != nil {
		out += " " + rs.ReturnValue.String()
	}
	return out
}

// CallExpression represents a function call (e.g., Add(10, 20))
type CallExpression struct {
	token.Token            // The '(' token
	Function    Expression // Identifier or expression evaluating to a function
	Arguments   []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Literal }
func (ce *CallExpression) String() string {
	var out string
	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out += ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
	return out
}
