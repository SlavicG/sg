package token

type TokenType string

const (
	// Identifiers + literals
	INT    = "INT"
	IDENT  = "IDENT"
	STRING = "STRING"

	// Operators
	SET    = "="
	PLUS   = "+"
	MINUS  = "-"
	EXC    = "!"
	STAR   = "*"
	SLASH  = "/"
	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA   = ","
	SEMICOL = ";"
	COL     = ":"
	LP      = "("
	RP      = ")"
	LB      = "{"
	RB      = "}"
	LBP     = "["
	RBP     = "]"

	//Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	//Miscellanios types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fun":      FUNCTION,
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"factos":   TRUE,
	"false":    FALSE,
	"unfactos": FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"ret":      RETURN,
}

func FindIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
