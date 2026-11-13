package token

type TokenType string

const (
	AS       = "AS"
	ASSIGN   = "="
	ASTERISK = "*"
	COMMA    = ","
	ELSE     = "ELSE"
	END      = "END"
	EOF      = "EOF"
	EQ       = "=="
	FALSE    = "FALSE"
	FOR      = "FOR"
	FUNCTION = "FUNCTION"
	GT       = ">"
	GTE      = ">="
	IDENT    = "IDENT"
	IF       = "IF"
	ILLEGAL  = "ILLEGAL"
	LBRACKET = "["
	LPAREN   = "("
	LT       = "<"
	LTE      = "<="
	MINUS    = "-"
	NONE     = "NONE"
	NOTEQ    = "!="
	NUMBER   = "NUMBER"
	NUMERIC  = "NUMERIC"
	PLUS     = "+"
	PRINT    = "PRINT"
	RBRACKET = "]"
	RETURN   = "RETURN"
	RPAREN   = ")"
	SLASH    = "/"
	STEP     = "STEP"
	STRING   = "STRING"
	THEN     = "THEN"
	TO       = "TO"
	TRUE     = "TRUE"
	VAR      = "VAR"
)

var keywords = map[string]TokenType{
	"As":       AS,
	"Else":     ELSE,
	"End":      END,
	"False":    FALSE,
	"For":      FOR,
	"Function": FUNCTION,
	"If":       IF,
	"None":     NONE,
	"Numeric":  NUMERIC,
	"Print":    PRINT,
	"Return":   RETURN,
	"Step":     STEP,
	"Then":     THEN,
	"To":       TO,
	"True":     TRUE,
	"Var":      VAR,
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
