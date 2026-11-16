package token

import "strings"

type TokenType string

const (
	AND      = "AND"
	AS       = "AS"
	ASSIGN   = "="
	ASTERISK = "*"
	CASE     = "CASE"
	COMMA    = ","
	CONST    = "CONST"
	ELSE     = "ELSE"
	ELSEIF   = "ELSEIF"
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
	OR       = "OR"
	PERIOD   = "."
	PLUS     = "+"
	RBRACKET = "]"
	RETURN   = "RETURN"
	RPAREN   = ")"
	SELECT   = "SELECT"
	SLASH    = "/"
	STEP     = "STEP"
	STRING   = "STRING"
	STRUCT   = "STRUCT"
	THEN     = "THEN"
	TO       = "TO"
	TRUE     = "TRUE"
	VAR      = "VAR"
)

var keywords = map[string]TokenType{
	"and":      AND,
	"as":       AS,
	"case":     CASE,
	"const":    CONST,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"end":      END,
	"false":    FALSE,
	"for":      FOR,
	"function": FUNCTION,
	"if":       IF,
	"none":     NONE,
	"numeric":  NUMERIC,
	"or":       OR,
	"return":   RETURN,
	"select":   SELECT,
	"step":     STEP,
	"struct":   STRUCT,
	"then":     THEN,
	"to":       TO,
	"true":     TRUE,
	"var":      VAR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
