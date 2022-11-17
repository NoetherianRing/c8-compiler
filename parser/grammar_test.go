package parser

import (
	"github.com/NoetherianRing/c8-compiler/token"
	"testing"
)

func TestBuild(t *testing.T){

	src := make([]token.Token, 0)
	/*slice := []token.Token{
		token.NewToken(token.LPAREN, token.LPAREN, 0),
		token.NewToken(token.INT, "3", 0),
	token.NewToken(token.PLUS, token.PLUS, 0),
	 token.NewToken(token.INT, "2", 0)
	}
	*/
	src = append(src, token.NewToken(token.LPAREN, token.LPAREN, 0))
	src = append(src, token.NewToken(token.INT, "3", 0))
	src = append(src, token.NewToken(token.PLUS, token.PLUS, 0))
	src = append(src, token.NewToken(token.INT, "2", 0))
	src = append(src, token.NewToken(token.RPAREN, token.RPAREN, 0))
	src = append(src, token.NewToken(token.ASTERISK, token.ASTERISK, 0))
	src = append(src, token.NewToken(token.INT, "80", 0))
	src = append(src, token.NewToken(token.EOF, token.EOF, 0))

//	index := 0

	grammar := GetGrammar()
	tree := NewSyntaxTree(NewNode(token.NewToken("","",0)))


//	valid := grammar["statements"].Build(src, &index)
	valid := grammar["statements"].Build(&src, tree)

	if valid != true || len(src) != 0{
		t.Fatalf("CHECK FALLO: (valid: %t len: %d)", valid, len(src))

	}else{
		tree.debug()
	}
}