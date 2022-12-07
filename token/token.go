package token

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	BYTE  = "BYTE"
	BOOL  = "BOOL"

	EQ   = "="
	EQEQ   = "=="
	AND   = "&"
	LAND   = "&&"
	OR   = "|"
	LOR   = "||"
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
	LBRACKET = "["
	RBRACKET = "]"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FN"
	LET = "LET"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"

	TYPEBOOL = "TYPEBOOL"
	TYPEBYTE = "TYPEBYTE"
	VOID = "VOID"

)

type Type string

type Token struct{
	Type    Type
	Literal string
	Line    int
}


func NewToken(tokenType Type, literal string, line int) Token{
	t := Token{Type: tokenType, Literal:literal, Line: line}
	return t
}
var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"bool":   TYPEBOOL,
	"byte":   TYPEBYTE,
	"true":   BOOL,
	"false":   BOOL,
	"void":   VOID,
}
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
