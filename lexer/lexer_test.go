package lexer

import(
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

)

func TestNextToken(t *testing.T){
	input, err := filepath.Abs("../fixtures/TestNextToken.txt")
	assert.NoError(t, err, "error in absPath TestNextToken")

	tests := []struct{
		expectedType token.Type
		expectedLiteral string
		expectedLine int

	}{
		{token.LET, "let",1},
		{token.IDENT, "five",1},
		{token.EQ, "=",1},
		{token.BYTE, "5",1},
		{token.SEMICOLON, ";",1},
		{token.NEWLINE, "\n",1},
		{token.LET, "let",2},
		{token.IDENT, "ten",2},
		{token.EQ, "=",2},
		{token.BYTE, "10",2},
		{token.SEMICOLON, ";",2},
		{token.NEWLINE, "\n",2},
		{token.LET, "let",3},
		{token.IDENT, "add",3},
		{token.EQ, "=",3},
		{token.FUNCTION, "fn",3},
		{token.LPAREN, "(",3},
		{token.IDENT, "x",3},
		{token.COMMA, ",",3},
		{token.IDENT, "y",3},
		{token.RPAREN, ")",3},
		{token.LBRACE, "{",3},
		{token.NEWLINE, "\n",3},
		{token.IDENT, "x",4},
		{token.PLUS, "+",4},
		{token.IDENT, "y",4},
		{token.SEMICOLON, ";",4},
		{token.NEWLINE, "\n",4},
		{token.NEWLINE, "\n",5},
		{token.RBRACE, "}", 6},
		{token.SEMICOLON, ";",6},
		{token.NEWLINE, "\n",6},
		{token.NEWLINE, "\n",7},
		{token.LET, "let",8},
		{token.IDENT, "result",8},
		{token.EQ, "=",8},
		{token.IDENT, "add",8},
		{token.LPAREN, "(",8},
		{token.IDENT, "five",8},
		{token.COMMA, ",",8},
		{token.IDENT, "ten",8},
		{token.RPAREN, ")",8},
		{token.SEMICOLON, ";",8},
		{token.NEWLINE, "\n",8},
		{token.BANG, "!",9},
		{token.MINUS, "-",9},
		{token.SLASH, "/",9},
		{token.ASTERISK, "*",9},
		{token.BYTE, "5",9},
		{token.SEMICOLON, ";",9},
		{token.NEWLINE, "\n",9},
		{token.BYTE, "5",10},
		{token.LT, "<",10},
		{token.BYTE, "10",10},
		{token.GT, ">",10},
		{token.BYTE, "5",10},
		{token.SEMICOLON, ";",10},
		{token.NEWLINE, "\n",10},
		{token.NEWLINE, "\n",11},
		{token.IF, "if",12},
		{token.LPAREN, "(",12},
		{token.BYTE, "5",12},
		{token.LT, "<",12},
		{token.BYTE, "10",12},
		{token.RPAREN, ")",12},
		{token.LBRACE, "{",12},
		{token.NEWLINE, "\n",12},
		{token.RETURN, "return",13},
		{token.BYTE, "1",13},
		{token.SEMICOLON, ";",13},
		{token.NEWLINE, "\n",13},
		{token.RBRACE, "}",14},
		{token.ELSE, "else",14},
		{token.LBRACE, "{",14},
		{token.NEWLINE, "\n",14},
		{token.RETURN, "return",15},
		{token.BYTE, "0",15},
		{token.SEMICOLON, ";",15},
		{token.NEWLINE, "\n",15},
		{token.RBRACE, "}",16},
		{token.NEWLINE, "\n",16},
		{token.NEWLINE, "\n",17},
		{token.BYTE, "10",18},
		{token.EQEQ, "==",18},
		{token.BYTE, "10",18},
		{token.SEMICOLON, ";",18},
		{token.NEWLINE, "\n",18},
		{token.BYTE, "10",19},
		{token.NOTEQ, "!=",19},
		{token.BYTE, "9",19},
		{token.SEMICOLON, ";",19},


		{token.EOF, "",19},
	}
	l, err := NewLexer(input)
	assert.NoError(t, err, "error in NewLexer TestNextToken")


	for i, tt := range tests{
		tok := l.nextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected =%q, got=%q",
				i, tt.expectedType, tok.Type)


		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected =%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)

		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected =%d, got=%d",
				i, tt.expectedLine, tok.Line)

		}
	}
}
