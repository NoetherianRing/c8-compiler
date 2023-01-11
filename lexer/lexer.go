package lexer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
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

	l := &Lexer{input: string(source), index: -1, cChar: "", cLine: 0}
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
			errorString := errorhandler.IllegalToken(t.Line, t.Literal)
			return nil, errors.New(errorString)
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
	case token.EQ:
		if l.peekChar() == token.EQ{
			l.readChar()
			tok = token.NewToken(token.EQEQ, token.EQEQ, l.cLine)
		}else{
			tok = token.NewToken(token.EQ, token.EQ, l.cLine)

		}
	case token.AND:
		if l.peekChar() == token.AND{
			l.readChar()
			tok = token.NewToken(token.LAND, token.LAND, l.cLine)
		}else{
			tok = token.NewToken(token.AND, token.AND, l.cLine)

		}
	case token.OR:
		if l.peekChar() == token.OR{
			l.readChar()
			tok = token.NewToken(token.LOR, token.LOR, l.cLine)
		}else{
			tok = token.NewToken(token.OR, token.OR, l.cLine)

		}
	case token.BANG:
		if l.peekChar() == token.EQ{
			l.readChar()
			tok = token.NewToken(token.NOTEQ, token.NOTEQ, l.cLine)
		}else {
			tok = token.NewToken(token.BANG, token.BANG, l.cLine)
		}
	case token.PLUS:
		tok = token.NewToken(token.PLUS, token.PLUS, l.cLine)
	case token.ASTERISK:
		tok = token.NewToken(token.ASTERISK, token.ASTERISK, l.cLine)
	case token.MINUS:
		tok = token.NewToken(token.MINUS, token.MINUS, l.cLine)
	case token.SLASH:
		tok = token.NewToken(token.SLASH, token.SLASH, l.cLine)
	case token.PERCENT:
		tok = token.NewToken(token.PERCENT, token.PERCENT, l.cLine)
	case token.DOLLAR:
		tok = token.NewToken(token.DOLLAR, token.DOLLAR, l.cLine)
	case token.GT:
		peek := l.peekChar()
		if peek == token.EQ{
			l.readChar()
			tok = token.NewToken(token.GTEQ, token.GTEQ, l.cLine)
		}else {
			if peek == token.GT{
				l.readChar()
				tok = token.NewToken(token.GTGT, token.GTGT, l.cLine)

			}else{
				tok = token.NewToken(token.GT, token.GT, l.cLine)

			}
		}
	case token.COMMA:
		tok = token.NewToken(token.COMMA, token.COMMA, l.cLine)
	case token.LPAREN:
		tok = token.NewToken(token.LPAREN, token.LPAREN, l.cLine)
	case token.RPAREN:
		tok = token.NewToken(token.RPAREN, token.RPAREN, l.cLine)
	case token.LBRACKET:
		tok = token.NewToken(token.LBRACKET, token.LBRACKET, l.cLine)
	case token.RBRACKET:
		tok = token.NewToken(token.RBRACKET, token.RBRACKET, l.cLine)
	case token.LBRACE:
		tok = token.NewToken(token.LBRACE, token.LBRACE, l.cLine)
	case token.RBRACE:
		tok = token.NewToken(token.RBRACE, token.RBRACE, l.cLine)
	case token.LT:
		peek := l.peekChar()
		if peek == token.EQ{
			l.readChar()
			tok = token.NewToken(token.LTEQ, token.LTEQ, l.cLine)
		}else {
			if peek == token.LT{
				l.readChar()
				tok = token.NewToken(token.LTLT, token.LTLT, l.cLine)

			}else{
				tok = token.NewToken(token.LT, token.LT, l.cLine)

			}
		}
	case token.XOR:
		tok = token.NewToken(token.XOR, token.XOR, l.cLine)
	case token.NEWLINE:
		tok = token.NewToken(token.NEWLINE, token.NEWLINE, l.cLine)
		l.cLine += 1
	case "":
		tok = token.NewToken(token.EOF, token.EOF, l.cLine)
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
	for isLetter(l.cChar) || isDigit(l.cChar){
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