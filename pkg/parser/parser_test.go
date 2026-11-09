package parser_test

import (
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