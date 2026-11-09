package parser

type TokenType string

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Identifiers & Literals
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Operators & Delimiters
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	LPAREN   = "("
	RPAREN   = ")"

	// Keywords
	VAR     = "VAR"
	AS      = "AS"
	NUMERIC = "NUMERIC"
	PRINT   = "PRINT"
)

var keywords = map[string]TokenType{
	"Var":     VAR,
	"As":      AS,
	"Numeric": NUMERIC,
	"Print":   PRINT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
