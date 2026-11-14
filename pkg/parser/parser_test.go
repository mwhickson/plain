package parser_test

import (
	"plain/pkg/ast"
	"plain/pkg/parser"
	"testing"
)

func TestVarStatement(t *testing.T) {
	input := `Var total As Numeric = 10 + 20`

	l := parser.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			t.Errorf("parser error: %v", err)
		}
		t.Fatalf("parser had %d errors", len(p.Errors()))
	}

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt := program.Statements[0]
	if stmt.String() != "Var total As Numeric = (10 + 20)" {
		t.Errorf("stmt.String() wrong. got=%q", stmt.String())
	}
}

func TestConstStatement(t *testing.T) {
	input := `Const answer As Numeric = 42`

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	if got := program.Statements[0].String(); got != "Const answer As Numeric = 42" {
		t.Errorf("statement string = %q, want %q", got, "Const answer As Numeric = 42")
	}
}

func TestVarStatementWithoutInitializer(t *testing.T) {
	input := "Var room As Room\nroom.Id = 100"

	p := parser.NewParser(parser.NewLexer(input))
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}
	if len(program.Statements) != 2 {
		t.Fatalf("statement count = %d, want 2", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("first statement has type %T, want *ast.VarStatement", program.Statements[0])
	}
	if stmt.Value != nil {
		t.Fatalf("initializer = %v, want nil", stmt.Value)
	}
}

func TestListLiteralAllowsTrailingComma(t *testing.T) {
	p := parser.NewParser(parser.NewLexer(`[1, 2,]`))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement has type %T, want *ast.ExpressionStatement", program.Statements[0])
	}
	list, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("expression has type %T, want *ast.ListLiteral", stmt.Expression)
	}
	if len(list.Elements) != 2 {
		t.Fatalf("list has %d elements, want 2", len(list.Elements))
	}
}

func TestCallArgumentsStillRejectTrailingComma(t *testing.T) {
	p := parser.NewParser(parser.NewLexer(`PrintLine("value",)`))
	p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Fatal("expected parser error for trailing comma in call arguments")
	}
}
