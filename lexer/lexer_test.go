package lexer

import(
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"

)

func TestNextToken(t *testing.T){

	type testsCases struct {
		description string
		fixture     string
		expectedTokens []token.Token
	}

	cases := []testsCases{
		{
			description: "TestNextToken1",
			fixture: "../fixtures/TestNextToken1.txt",
			expectedTokens: []token.Token{

				token.NewToken(token.LBRACE, token.LBRACE, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.LET, "let", 1),
				token.NewToken(token.IDENT, "foo", 1),
				token.NewToken(token.TYPEBYTE, "byte", 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 2),
				token.NewToken(token.LET, "let", 3),
				token.NewToken(token.IDENT, "_f00", 3),
				token.NewToken(token.ASTERISK, token.ASTERISK, 3),
				token.NewToken(token.ASTERISK, token.ASTERISK, 3),
				token.NewToken(token.LBRACKET, token.LBRACKET, 3),
				token.NewToken(token.BYTE, "30", 3),
				token.NewToken(token.RBRACKET, token.RBRACKET, 3),
				token.NewToken(token.TYPEBYTE, "byte", 3),
				token.NewToken(token.NEWLINE, token.NEWLINE, 3),
				token.NewToken(token.LET, "let", 4),
				token.NewToken(token.IDENT, "f_o0", 4),
				token.NewToken(token.MINUS, token.MINUS, 4),
				token.NewToken(token.ILLEGAL, ";", 4),
				token.NewToken(token.NEWLINE, token.NEWLINE, 4),
				token.NewToken(token.RBRACE, token.RBRACE, 5),
				token.NewToken(token.EOF, token.EOF, 5),

			},

		},
		{
			description: "TestNextToken2",
			fixture: "../fixtures/TestNextToken2.txt",
			expectedTokens: []token.Token{
				token.NewToken(token.WHILE, "while", 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.DOLLAR, token.DOLLAR, 0),
				token.NewToken(token.IDENT, "var", 0),
				token.NewToken(token.NOTEQ, token.NOTEQ, 0),
				token.NewToken(token.LBRACKET, token.LBRACKET, 0),
				token.NewToken(token.IDENT, "i", 0),
				token.NewToken(token.RBRACKET, token.RBRACKET, 0),
				token.NewToken(token.IDENT, "address", 0),
				token.NewToken(token.LTLT, token.LTLT, 0),
				token.NewToken(token.BYTE, "10", 0),
				token.NewToken(token.GTGT, token.GTGT, 0),
				token.NewToken(token.BYTE, "20", 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.LAND, token.LAND, 0),
				token.NewToken(token.IDENT, "var2", 0),
				token.NewToken(token.OR, token.OR, 0),
				token.NewToken(token.IDENT, "var3", 0),
				token.NewToken(token.AND, token.AND, 0),
				token.NewToken(token.IDENT, "var4", 0),
				token.NewToken(token.XOR, token.XOR, 0),
				token.NewToken(token.IDENT, "var5", 0),
				token.NewToken(token.EQEQ, token.EQEQ, 0),
				token.NewToken(token.IDENT, "call", 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.LOR, token.LOR, 0),
				token.NewToken(token.BYTE, "5", 0),
				token.NewToken(token.SLASH, token.SLASH, 0),
				token.NewToken(token.BYTE, "2", 0),
				token.NewToken(token.LTEQ, token.LTEQ, 0),
				token.NewToken(token.BYTE, "2", 0),
				token.NewToken(token.PERCENT, token.PERCENT, 0),
				token.NewToken(token.BYTE, "3", 0),
				token.NewToken(token.LT, token.LT, 0),
				token.NewToken(token.MINUS, token.MINUS, 0),
				token.NewToken(token.BYTE, "1", 0),
				token.NewToken(token.GT, token.GT, 0),
				token.NewToken(token.PLUS, token.PLUS, 0),
				token.NewToken(token.BYTE, "2", 0),
				token.NewToken(token.GTEQ, token.GTEQ, 0),
				token.NewToken(token.BYTE, "5", 0),
				token.NewToken(token.EOF, token.EOF, 0),

			},

		},
		{
			description: "TestNextToken3",
			fixture: "../fixtures/TestNextToken3.txt",
			expectedTokens: []token.Token{
				token.NewToken(token.IDENT, "foo", 0),
				token.NewToken(token.EQ, token.EQ, 0),
				token.NewToken(token.IDENT, "call", 0),
				token.NewToken(token.LPAREN, token.LPAREN, 0),
				token.NewToken(token.BOOL, "true", 0),
				token.NewToken(token.COMMA, token.COMMA, 0),
				token.NewToken(token.BOOL, "false", 0),
				token.NewToken(token.RPAREN, token.RPAREN, 0),
				token.NewToken(token.NEWLINE, token.NEWLINE, 0),
				token.NewToken(token.IF, "if", 1),
				token.NewToken(token.BOOL, "true", 1),
				token.NewToken(token.LBRACE, token.LBRACE, 1),
				token.NewToken(token.NEWLINE, token.NEWLINE, 1),
				token.NewToken(token.FUNCTION, "fn", 2),
				token.NewToken(token.LPAREN, token.LPAREN, 2),
				token.NewToken(token.LET, "let", 2),
				token.NewToken(token.IDENT, "arg", 2),
				token.NewToken(token.TYPEBOOL, "bool", 2),
				token.NewToken(token.RPAREN, token.RPAREN, 2),
				token.NewToken(token.IDENT, "function", 2),
				token.NewToken(token.VOID, "void", 2),
				token.NewToken(token.LBRACE, token.LBRACE, 2),
				token.NewToken(token.NEWLINE, token.NEWLINE, 2),
				token.NewToken(token.RETURN, "return", 3),
				token.NewToken(token.BOOL, "false", 3),
				token.NewToken(token.NEWLINE, token.NEWLINE, 3),
				token.NewToken(token.RBRACE, token.RBRACE, 4),
				token.NewToken(token.NEWLINE, token.NEWLINE, 4),
				token.NewToken(token.RBRACE, token.RBRACE, 5),
				token.NewToken(token.ELSE, "else", 5),
				token.NewToken(token.LBRACE, token.LBRACE, 5),
				token.NewToken(token.NEWLINE, token.NEWLINE, 5),
				token.NewToken(token.IDENT, "foo", 6),
				token.NewToken(token.EQ, token.EQ, 6),
				token.NewToken(token.BANG, token.BANG, 6),
				token.NewToken(token.IDENT, "foo", 6),
				token.NewToken(token.NEWLINE, token.NEWLINE, 6),
				token.NewToken(token.RBRACE, token.RBRACE, 7),
				token.NewToken(token.EOF, token.EOF, 7),












			},
		},
	}
	for i, tt := range cases{
		input, err := filepath.Abs(tt.fixture)
		assert.NoError(t, err, "error in absPath  " + tt.fixture )

		l, err := NewLexer(input)
		assert.NoError(t, err, "error in NewLexer " + tt.description)
		for j, ttt := range tt.expectedTokens{
			tok := l.nextToken()


			if tok.Type != ttt.Type {
				t.Fatalf("tests[%d] - token[%d]type wrong. expected =%q, got=%q",
					i, j, ttt.Type, tok.Type)


			}

			if tok.Literal != ttt.Literal {
				t.Fatalf("tests[%d] - token[%d]literal wrong. expected =%q, got=%q",
					i,j, ttt.Literal, tok.Literal)

			}

			if tok.Line != ttt.Line{
				t.Fatalf("tests[%d] - token[%d]line wrong. expected =%d, got=%d",
					i,j, ttt.Line, tok.Line)

			}
		}

	}
}
