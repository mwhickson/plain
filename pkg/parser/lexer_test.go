package parser_test

import (
	"plain/pkg/parser"
	"testing"
)

func TestLexerBasic(t *testing.T) {
	input := `Var total As Numeric = 10 + 20`

	tests := []struct {
		expectedType    parser.TokenType
		expectedLiteral string
	}{
		{parser.VAR, "Var"},
		{parser.IDENT, "total"},
		{parser.AS, "As"},
		{parser.NUMERIC, "Numeric"},
		{parser.ASSIGN, "="},
		{parser.NUMBER, "10"},
		{parser.PLUS, "+"},
		{parser.NUMBER, "20"},
		{parser.EOF, ""},
	}

	l := parser.NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}