package parser_test

import (
	"plain/pkg/parser"
	"plain/pkg/token"
	"testing"
)

func TestLexerBasic(t *testing.T) {
	input := `Var total As Numeric = 10 + 20`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "Var"},
		{token.IDENT, "total"},
		{token.AS, "As"},
		{token.NUMERIC, "Numeric"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.PLUS, "+"},
		{token.NUMBER, "20"},
		{token.EOF, ""},
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
