package token

type TokenType string

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Identifiers & Literals
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// Operators & Delimiters
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	EQ       = "=="
	LPAREN   = "("
	RPAREN   = ")"
	COMMA    = ","
	NOTEQ    = "!="
	LT       = "<"
	LTE      = "<="
	GT       = ">"
	GTE      = ">="

	// Keywords
	VAR      = "VAR"
	AS       = "AS"
	NUMERIC  = "NUMERIC"
	PRINT    = "PRINT"
	IF       = "IF"
	THEN     = "THEN"
	ELSE     = "ELSE"
	END      = "END"
	FOR      = "FOR"
	TO       = "TO"
	STEP     = "STEP"
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	NONE     = "NONE"
)

var keywords = map[string]TokenType{
	"Var":      VAR,
	"As":       AS,
	"Numeric":  NUMERIC,
	"Print":    PRINT,
	"If":       IF,
	"Then":     THEN,
	"Else":     ELSE,
	"End":      END,
	"For":      FOR,
	"To":       TO,
	"Step":     STEP,
	"Function": FUNCTION,
	"Return":   RETURN,
	"None":     NONE,
	"True":     TRUE,
	"False":    FALSE,
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
