package ast

import (
	"plain/pkg/token"
	"strings"
)

//
// AssignmentStatement
//

type AssignmentStatement struct {
	token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Literal }
func (as *AssignmentStatement) String() string {
	return as.Name.String() + " = " + as.Value.String()
}

//
// BooleanLiteral
//

type BooleanLiteral struct {
	token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Literal }

//
// CallExpression represents a function call (e.g., Add(10, 20))
//

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

//
// Expression
//

type Expression interface {
	Node
	expressionNode()
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

//
// ForStatement
//

type ForStatement struct {
	token.Token
	VarName  *Identifier
	TypeName string
	Start    Expression
	End      Expression
	Step     Expression // Optional; can be nil (defaults to 1)
	Body     []Statement
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

//
// FunctionStatement
//

type FunctionStatement struct {
	token.Token
	Name       *Identifier
	Parameters []*Parameter
	ReturnType string // Optional; empty if none
	Body       []Statement
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

//
// Identifier
//

type Identifier struct {
	token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Literal }
func (i *Identifier) String() string       { return i.Value }

//
// IfStatement
//

type IfStatement struct {
	token.Token
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

//
// IndexExpression represents an index access like items[0]
//

type IndexExpression struct {
	token.Token // The '[' token
	Left        Expression
	Index       Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

//
// InfixExpression
//

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

//
// IntegerLiteral
//

type IntegerLiteral struct {
	token.Token
	Value float64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Literal }
func (il *IntegerLiteral) String() string       { return il.Literal }

//
// ListLiteral
//

type ListLiteral struct {
	token.Token // The '[' token
	Elements    []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Literal }
func (ll *ListLiteral) String() string {
	var out string
	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}
	out += "[" + strings.Join(elements, ", ") + "]"
	return out
}

//
// Node
//

type Node interface {
	TokenLiteral() string
	String() string
}

//
// NoneLiteral
//

type NoneLiteral struct {
	token.Token
}

func (nl *NoneLiteral) expressionNode()      {}
func (nl *NoneLiteral) TokenLiteral() string { return nl.Literal }
func (nl *NoneLiteral) String() string       { return "None" }

//
// Parameter
//

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

//
// Program
//

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

//
// PropertyAccessExpression represents dot notation like v.X
//

type PropertyAccessExpression struct {
	token.Token // The '.' token
	Left        Expression
	Property    *Identifier
}

func (pae *PropertyAccessExpression) expressionNode()      {}
func (pae *PropertyAccessExpression) TokenLiteral() string { return pae.Literal }
func (pae *PropertyAccessExpression) String() string {
	return "(" + pae.Left.String() + "." + pae.Property.String() + ")"
}

//
// ReturnStatement
//

type ReturnStatement struct {
	token.Token
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

//
// Statement
//

type Statement interface {
	Node
	statementNode()
}

//
// StringLiteral
//

type StringLiteral struct {
	token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Literal }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

//
// StructMember
//

type StructMember struct {
	token.Token
	Name     *Identifier
	TypeName string
}

func (sm *StructMember) statementNode()       {}
func (sm *StructMember) TokenLiteral() string { return sm.Literal }
func (sm *StructMember) String() string {
	return sm.Name.String() + " As " + sm.TypeName
}

//
// StructStatement
//

type StructStatement struct {
	token.Token
	Name    *Identifier
	Members []*StructMember
}

func (ss *StructStatement) statementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Literal }
func (ss *StructStatement) String() string {
	var out string
	out += "Struct " + ss.Name.String() + " "
	members := []string{}
	for _, m := range ss.Members {
		members = append(members, m.String())
	}
	out += strings.Join(members, " ") + " End Struct"
	return out
}

//
// VarStatement
//

type VarStatement struct {
	token.Token
	Name     *Identifier
	TypeName string
	Value    Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Literal }
func (vs *VarStatement) String() string {
	return "Var " + vs.Name.String() + " As " + vs.TypeName + " = " + vs.Value.String()
}
