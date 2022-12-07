package lexer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/token"
	"os"
)

type Lexer struct{
	input string
	index int    //current position
	cChar string //current char
	cLine int    //current line
}

func NewLexer(filename string) (*Lexer, error){

	source, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	l := &Lexer{input: string(source), index: 0, cChar: "", cLine: 1}
	l.readChar()
	return l, nil
}

func (l *Lexer) readChar(){
	if l.index+1 >= len(l.input){
		l.cChar = ""
	}else{
		l.cChar = string(l.input[l.index+1 ])
	}
	l.index = l.index +1
}
func (l *Lexer) GetTokens() ([]token.Token, error){
	tokens := make([]token.Token, 0)
	t := l.nextToken()
	for t.Type != token.EOF{
		if t.Type == token.ILLEGAL{
			return nil, errors.New("illegal token: " + t.Literal + " in line: " + string(rune(t.Line)))
		}
		tokens = append(tokens, t)
		t = l.nextToken()
	}
	tokens = append(tokens, t)
	return tokens, nil
}
func (l *Lexer) nextToken() token.Token{
	var tok token.Token

	l.skipWhitespace()
	l.skipComment()

	switch l.cChar{
	case "=":
		if l.peekChar() == "="{
			ch := l.cChar
			l.readChar()
			tok = token.NewToken(token.EQEQ, ch + l.cChar, l.cLine)
		}else{
			tok = token.NewToken(token.EQ, l.cChar, l.cLine)

		}
	case "&":
		if l.peekChar() == "&"{
			ch := l.cChar
			l.readChar()
			tok = token.NewToken(token.LAND, ch + l.cChar, l.cLine)
		}else{
			tok = token.NewToken(token.AND, l.cChar, l.cLine)

		}

	case "|":
		if l.peekChar() == "|"{
			ch := l.cChar
			l.readChar()
			tok = token.NewToken(token.LOR, ch + l.cChar, l.cLine)
		}else{
			tok = token.NewToken(token.OR, l.cChar, l.cLine)

		}
	case "!":
		if l.peekChar() == "="{
			ch := l.cChar
			l.readChar()
			tok = token.NewToken(token.NOTEQ, ch + l.cChar, l.cLine)
		}else {
			tok = token.NewToken(token.BANG, l.cChar, l.cLine)
		}
	case "+":
		tok = token.NewToken(token.PLUS, l.cChar, l.cLine)
	case "*":
		tok = token.NewToken(token.ASTERISK, l.cChar, l.cLine)
	case "-":
		tok = token.NewToken(token.MINUS, l.cChar, l.cLine)
	case "/":
		tok = token.NewToken(token.SLASH, l.cChar, l.cLine)
	case "%":
		tok = token.NewToken(token.PERCENT, l.cChar, l.cLine)
	case ">":
		tok = token.NewToken(token.GT, l.cChar, l.cLine)
	case ",":
		tok = token.NewToken(token.COMMA, l.cChar, l.cLine)
	case "(":
		tok = token.NewToken(token.LPAREN, l.cChar, l.cLine)
	case ")":
		tok = token.NewToken(token.RPAREN, l.cChar, l.cLine)
	case "[":
		tok = token.NewToken(token.LBRACKET, l.cChar, l.cLine)
	case "]":
		tok = token.NewToken(token.RBRACKET, l.cChar, l.cLine)
	case "{":
		tok = token.NewToken(token.LBRACE, l.cChar, l.cLine)
	case "}":
		tok = token.NewToken(token.RBRACE, l.cChar, l.cLine)
	case "<":
		tok = token.NewToken(token.LT, l.cChar, l.cLine)
	case ";":
		tok = token.NewToken(token.SEMICOLON, l.cChar, l.cLine)
	case "\n":
		tok = token.NewToken(token.NEWLINE, l.cChar, l.cLine)
		l.cLine += 1
	case "":
		tok = token.NewToken(token.EOF, l.cChar, l.cLine)
	default:
		if isLetter(l.cChar){
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.cLine
			return tok
		}else if isDigit(l.cChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.BYTE
			tok.Line = l.cLine
			return tok
		}else{
			tok = token.NewToken(token.ILLEGAL, l.cChar, l.cLine)
		}
	}
	l.readChar()
	return tok

}
func (l *Lexer) readNumber() string {
	position := l.index
	for isDigit(l.cChar) {
		l.readChar()
	}
	return l.input[position:l.index]
}
func isDigit(ch string) bool {
	return "0" <= ch && ch <= "9"
}
func (l *Lexer) readIdentifier() string {
	startPosition := l.index
	for isLetter(l.cChar) {
		l.readChar()
	}
	return l.input[startPosition:l.index]
}
func isLetter(ch string) bool {
	return "a" <= ch && ch <= "z" || "A" <= ch && ch <= "Z" || ch == "_"

}

func (l *Lexer) skipWhitespace(){
	for l.cChar == " " || l.cChar == "\t" || l.cChar == "\r"{
		l.readChar()
	}
}


func (l *Lexer) skipComment(){
	if l.cChar == "#"{
		for l.cChar != "\n"{
			l.readChar()
		}
	}
}

func (l *Lexer) peekChar() string{
	return string(l.input[l.index+1])
}