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

func TestLexerDecimalAndPropertyAccessNumbers(t *testing.T) {
	l := parser.NewLexer("12.5 12.field")
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NUMBER, "12.5"},
		{token.NUMBER, "12"},
		{token.PERIOD, "."},
		{token.IDENT, "field"},
		{token.EOF, ""},
	}
	for i, expected := range tests {
		got := l.NextToken()
		if got.Type != expected.expectedType || got.Literal != expected.expectedLiteral {
			t.Fatalf("token %d = (%q, %q), want (%q, %q)", i, got.Type, got.Literal, expected.expectedType, expected.expectedLiteral)
		}
	}
}

func TestLexerLineContinuation(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"LF", "first \\\n second"},
		{"CRLF", "first\\\r\nsecond"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := parser.NewLexer(tt.input)
			first := l.NextToken()
			second := l.NextToken()

			if first.Type != token.IDENT || first.Literal != "first" || first.Line != 1 {
				t.Fatalf("first token = %+v, want first identifier on line 1", first)
			}
			if second.Type != token.IDENT || second.Literal != "second" || second.Line != 2 {
				t.Fatalf("second token = %+v, want second identifier on line 2", second)
			}
			if tok := l.NextToken(); tok.Type != token.EOF {
				t.Fatalf("final token = %+v, want EOF", tok)
			}
		})
	}
}
