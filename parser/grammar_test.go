package parser

import (
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T){

	grammar := GetGrammar()

	type cases struct{
		description string
		src []token.Token
		isValid bool
		desiredTreeRep string
	}

	testCases := []cases{
		{
			description: "foo = (3+2)*80",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),

				token.NewToken(token.IDENT, "foo", 0),
				token.NewToken(token.EQ, token.EQ, 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.BYTE, "3", 0),
				token.NewToken(token.PLUS, token.PLUS, 0),
				token.NewToken(token.BYTE, "2", 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.ASTERISK, token.ASTERISK, 0),
				token.NewToken(token.BYTE, "80", 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),

				token.NewToken(token.RBRACE, token.RBRACE, 1),

				token.NewToken(token.EOF, token.EOF, 1),
			},
			isValid: true,
			//Tree expressed as an string:
			/* 	EOF
				  |
			      }
			      |
			      =
			    /   \
			   id    *
			        / \
			       )   80
			       |
			       +
			     /  \
			   3     2
			*/
			desiredTreeRep: "\n/EOF\n" +
				"/EOF/}\n" +
				"/EOF/}/=\n" +
				"/EOF/}/=/foo\n" +
				"/EOF/}/=/*\n" +
				"/EOF/}/=/*/)\n" +
				"/EOF/}/=/*/)/+\n" +
				"/EOF/}/=/*/)/+/3\n" +
				"/EOF/}/=/*/)/+/2\n" +
				"/EOF/}/=/*/80\n",
		},
		{
			description: "let *byte / 4 new lines",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),

				token.NewToken(token.LET, "let", 1),
				token.NewToken(token.IDENT, "foo", 1),
				token.NewToken(token.ASTERISK, "*", 1),
				token.NewToken(token.TYPEBYTE, "BYTE", 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),

				token.NewToken(token.NEWLINE, token.NEWLINE, 2),

				token.NewToken(token.NEWLINE, token.NEWLINE, 3),

				token.NewToken(token.RBRACE, token.RBRACE, 4),
				token.NewToken(token.EOF, token.EOF, 5),

			},
			isValid: true,
			desiredTreeRep: "\n/EOF\n" +
				"/EOF/}\n" +
				"/EOF/}/let\n" +
				"/EOF/}/let/foo\n" +
				"/EOF/}/let/*\n" +
				"/EOF/}/let/*/BYTE\n",

		},
		{
			description: "new line / let [30]bool",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),

				token.NewToken(token.LET, "let", 1),
				token.NewToken(token.IDENT, "foo", 1),
				token.NewToken(token.LBRACKET, "[", 1),
				token.NewToken(token.BYTE, "30", 1),
				token.NewToken(token.RBRACKET, "]", 1),
				token.NewToken(token.TYPEBOOL, "bool", 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),


				token.NewToken(token.RBRACE, token.RBRACE, 3),
				token.NewToken(token.EOF, token.EOF, 4),

			},
			isValid: true,
			desiredTreeRep: "\n/EOF\n" +
				"/EOF/}\n" +
				"/EOF/}/let\n" +
				"/EOF/}/let/foo\n" +
				"/EOF/}/let/]\n" +
				"/EOF/}/let/]/30\n" +
				"/EOF/}/let/]/bool\n",

		},
		{
			description: "call()",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.IDENT, "foo", 1),
				token.NewToken(token.LPAREN, "(", 1),
				token.NewToken(token.RPAREN, ")", 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.RBRACE, token.RBRACE, 2),
				token.NewToken(token.EOF, token.EOF, 2),
			},
			isValid: true,
			desiredTreeRep: "\n/EOF\n" +
				"/EOF/}\n" +
				"/EOF/}/)\n"+
				"/EOF/}/)/foo\n",

		},
		{
			description: "call(var1, &var2, 3, true, 7+8)",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.IDENT, "call", 1),
				token.NewToken(token.LPAREN, "(", 1),
				token.NewToken(token.IDENT, "var1", 1),
				token.NewToken(token.COMMA, ",", 1),
				token.NewToken(token.AND, "&", 1),
				token.NewToken(token.AND, "&", 1),
				token.NewToken(token.IDENT, "var2", 1),
				token.NewToken(token.COMMA, ",", 1),
				token.NewToken(token.BYTE, "3", 1),
				token.NewToken(token.COMMA, ",", 1),
				token.NewToken(token.BOOL, "true", 1),
				token.NewToken(token.COMMA, ",", 1),
				token.NewToken(token.BYTE, "7", 1),
				token.NewToken(token.PLUS, "+", 1),
				token.NewToken(token.BYTE, "8", 1),
				token.NewToken(token.RPAREN, ")", 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.RBRACE, token.RBRACE, 2),
				token.NewToken(token.EOF, token.EOF, 2),
			},
			isValid: true,
			desiredTreeRep:
				"\n/EOF\n" +
				"/EOF/}" +
				"\n/EOF/}/)" +
				"\n/EOF/}/)/call\n" +
				"/EOF/}/)/," +
				"\n/EOF/}/)/,/var1" +
				"\n/EOF/}/)/,/," +
				"\n/EOF/}/)/,/,/&" +
				"\n/EOF/}/)/,/,/&/&\n" +
				"/EOF/}/)/,/,/&/&/var2\n" +
				"/EOF/}/)/,/,/," +
				"\n/EOF/}/)/,/,/,/3" +
				"\n/EOF/}/)/,/,/,/,\n" +
				"/EOF/}/)/,/,/,/,/true\n" +
				"/EOF/}/)/,/,/,/,/+\n" +
				"/EOF/}/)/,/,/,/,/+/7" +
				"\n/EOF/}/)/,/,/,/,/+/8\n",
			},{
				description: "call(var1)",
				src: []token.Token{
					token.NewToken(token.LBRACE, token.LBRACE, 0),
					token.NewToken(token.NEWLINE, token.NEWLINE, 0),
					token.NewToken(token.IDENT, "call", 1),
					token.NewToken(token.LPAREN, "(", 1),
					token.NewToken(token.IDENT, "var1", 1),
					token.NewToken(token.RPAREN, ")", 1),
					token.NewToken(token.NEWLINE, token.NEWLINE, 1),
					token.NewToken(token.RBRACE, token.RBRACE, 2),
					token.NewToken(token.EOF, token.EOF, 2),
				},
				isValid: true,
				desiredTreeRep:
					"\n/EOF\n" +
					"/EOF/}" +
					"\n/EOF/}/)" +
					"\n/EOF/}/)/call\n" +
					"/EOF/}/)/var1\n",
			},
			{
			  description: "**var1=*[2]([8][10]matrix)",
			  src: []token.Token{
				  token.NewToken(token.LBRACE, token.LBRACE, 0),
				  token.NewToken(token.ASTERISK, token.ASTERISK, 0),
				  token.NewToken(token.ASTERISK, token.ASTERISK, 0),
				  token.NewToken(token.IDENT, "var1", 0),
				  token.NewToken(token.EQ, token.EQ, 0),
				  token.NewToken(token.ASTERISK, token.ASTERISK, 0),
				  token.NewToken(token.LBRACKET, token.LBRACKET, 0),
				  token.NewToken(token.BYTE, "2", 0),
				  token.NewToken(token.RBRACKET, token.RBRACKET, 0),
				  token.NewToken(token.LPAREN, token.LPAREN, 0),
				  token.NewToken(token.LBRACKET, token.LBRACKET, 0),
				  token.NewToken(token.BYTE, "8", 0),
				  token.NewToken(token.RBRACKET, token.RBRACKET, 0),
				  token.NewToken(token.LBRACKET, token.LBRACKET, 0),
				  token.NewToken(token.BYTE, "10", 0),
				  token.NewToken(token.RBRACKET, token.RBRACKET, 0),
				  token.NewToken(token.IDENT, "matrix", 0),
				  token.NewToken(token.RPAREN, token.RPAREN, 0),
				  token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				  token.NewToken(token.RBRACE, token.RBRACE, 1),
				  token.NewToken(token.EOF, token.EOF, 1),

			  },
			  isValid: true,
			  desiredTreeRep: "\n/EOF\n" +
			  	"/EOF/}\n" +
			  	"/EOF/}/=\n" +
			  	"/EOF/}/=/*\n"+
			  	"/EOF/}/=/*/*\n"+
			  	"/EOF/}/=/*/*/var1\n"+
			  	"/EOF/}/=/*\n"+
			  	"/EOF/}/=/*/]\n"+
			  	"/EOF/}/=/*/]/2\n"+
			  	"/EOF/}/=/*/]/)\n"+
			  	"/EOF/}/=/*/]/)/]\n"+
			  	"/EOF/}/=/*/]/)/]/8\n"+
			  	"/EOF/}/=/*/]/)/]/]\n"+
			  	"/EOF/}/=/*/]/)/]/]/10\n"+
			  	"/EOF/}/=/*/]/)/]/]/matrix\n",
			},
			{
				description: "fn myFunc() void {new line}",
				src: []token.Token{
					token.NewToken(token.LBRACE, token.LBRACE, 0),
					token.NewToken(token.FUNCTION, "fn", 0),
					token.NewToken(token.IDENT, "myFunc", 0),
					token.NewToken(token.LPAREN, token.LPAREN, 0),
					token.NewToken(token.RPAREN, token.RPAREN, 0),
					token.NewToken(token.VOID, "void", 0),
					token.NewToken(token.LBRACE, token.LBRACE, 0),
					token.NewToken(token.NEWLINE, token.NEWLINE, 0),
					token.NewToken(token.RBRACE, token.RBRACE, 1),
					token.NewToken(token.NEWLINE, token.NEWLINE, 1),
					token.NewToken(token.RBRACE, token.RBRACE, 2),
					token.NewToken(token.EOF, token.EOF, 2),
				},
				isValid: true,
				desiredTreeRep:
					"\n/EOF\n" +
					"/EOF/}\n"+
					"/EOF/}/fn\n"+
					"/EOF/}/fn/myFunc\n"+
					"/EOF/}/fn/)\n"+
					"/EOF/}/fn/void\n"+
					"/EOF/}/fn/}\n",
			},
		{
			description: "fn myFunc(let number byte) byte {new line}",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.FUNCTION, "fn", 0),
				token.NewToken(token.IDENT, "myFunc", 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.LET, "let", 0),
				token.NewToken(token.IDENT, "number", 0),
				token.NewToken(token.TYPEBYTE, "byte", 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.TYPEBYTE, "byte", 0),
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.RBRACE, token.RBRACE, 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.RBRACE, token.RBRACE, 2),
				token.NewToken(token.EOF, token.EOF, 2),
			},
			isValid: true,
			desiredTreeRep:
			"\n/EOF\n" +
				"/EOF/}\n"+
				"/EOF/}/fn\n"+
				"/EOF/}/fn/myFunc\n"+
				"/EOF/}/fn/)\n"+
				"/EOF/}/fn/)/let\n"+
				"/EOF/}/fn/)/let/number\n"+
				"/EOF/}/fn/)/let/byte\n"+
				"/EOF/}/fn/byte\n"+
				"/EOF/}/fn/}\n",
		},
		{
			description: "fn myFunc(let number byte, let flag bool) byte {new line}",
			src: []token.Token{
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.FUNCTION, "fn", 0),
				token.NewToken(token.IDENT, "myFunc", 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.LET, "let", 0),
				token.NewToken(token.IDENT, "number", 0),
				token.NewToken(token.TYPEBYTE, "byte", 0),
				token.NewToken(token.COMMA, token.COMMA, 0),
				token.NewToken(token.LET, "let", 0),
				token.NewToken(token.IDENT, "flag", 0),
				token.NewToken(token.TYPEBOOL, "bool", 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.TYPEBYTE, "byte", 0),
				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.RBRACE, token.RBRACE, 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.RBRACE, token.RBRACE, 2),
				token.NewToken(token.EOF, token.EOF, 2),
			},
			isValid: true,
			desiredTreeRep:
			"\n/EOF\n" +
				"/EOF/}\n"+
				"/EOF/}/fn\n"+
				"/EOF/}/fn/myFunc\n"+
				"/EOF/}/fn/)\n"+
				"/EOF/}/fn/)/,\n"+
				"/EOF/}/fn/)/,/let\n"+
				"/EOF/}/fn/)/,/let/number\n"+
				"/EOF/}/fn/)/,/let/byte\n"+
				"/EOF/}/fn/)/,/let\n"+
				"/EOF/}/fn/)/,/let/flag\n"+
				"/EOF/}/fn/)/,/let/bool\n"+
				"/EOF/}/fn/byte\n"+
				"/EOF/}/fn/}\n",
		},
	}

	for _, scenario := range testCases{

		t.Run(scenario.description, func(t *testing.T) {
			tree := NewSyntaxTree(NewNode(token.NewToken("","",0)))
			valid := grammar["program"].Build(&scenario.src, tree)
			assert.Equal(t, scenario.isValid, valid)
			treeRep :=""
			parseTree(tree, "", &treeRep)
			assert.Equal(t, scenario.desiredTreeRep, treeRep)
		})
	}

}


func parseTree(tree *SyntaxTree, parents string, treeRep *string){
	*treeRep += parents + tree.current.value.Literal +"\n"
	for _, child := range tree.current.children{
		current := tree.current
		tree.current = child
		parseTree(tree, parents+ current.value.Literal +"/", treeRep)
		tree.current = current
	}

}