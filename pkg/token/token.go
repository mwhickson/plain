package token

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
	IF      = "IF"
	THEN    = "THEN"
	ELSE    = "ELSE"
	END     = "END"
	FOR     = "FOR"
	TO      = "TO"
	STEP    = "STEP"
)

var keywords = map[string]TokenType{
	"Var":     VAR,
	"As":      AS,
	"Numeric": NUMERIC,
	"Print":   PRINT,
	"If":      IF,
	"Then":    THEN,
	"Else":    ELSE,
	"End":     END,
	"For":     FOR,
	"To":      TO,
	"Step":    STEP,
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
