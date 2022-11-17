package token

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT  = "INT"

	EQ   = "="
	EQEQ   = "=="
	NOTEQ   = "!="
	PLUS = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"
	PERCENT = "%"

	LT = "<"
	GT = ">"

	COMMA = ","
	SEMICOLON = ";"
	NEWLINE = "\n"

	LPAREN = "("
	RPAREN = ")"
	RBRACKET = "["
	LBRACKET = "]"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FN"
	LET = "LET"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"

	TYPEBOOL = "TYPEBOOL"
	TYPEINT  = "TYPEINT"
	TYPEBYTE = "TYPEBYTE"

)

type TokenType string

type Token struct{
	Type    TokenType
	Literal string
	Line int
}


func NewToken(tokenType TokenType, literal string, line int) Token{
	t := Token{Type: tokenType, Literal:literal, Line: line}
	return t
}
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"int":    TYPEINT,
	"bool":   TYPEBOOL,
	"byte":   TYPEBYTE,
}
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
