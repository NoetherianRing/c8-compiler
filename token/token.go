package token

const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	BYTE  = "BYTE"
	BOOL  = "BOOL"

	TRUE  = "true"
	FALSE  = "false"

	EQ   = "="

	DOLLAR   = "$"

	AND   = "&"
	LAND   = "&&"
	OR       = "|"
	LOR      = "||"
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
	LTLT     = "<<"
	GTGT     = ">>"
	XOR      = "^"

	LT = "<"
	LTEQ = "<="
	GT = ">"
	GTEQ = ">="
	NOTEQ   = "!="
	EQEQ   = "=="

	COMMA = ","
	NEWLINE = "\n"

	LPAREN = "("
	RPAREN = ")"
	LBRACKET = "["
	RBRACKET = "]"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "fn"
	WHILE = "while"
	LET = "let"
	IF = "if"
	ELSE = "else"
	RETURN = "return"
	MAIN = "main"

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
	"while": WHILE,
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
