package parser

import (
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func GetTree(path string) *ast.SyntaxTree{
	absPath, err := filepath.Abs(path)
	if err != nil{
		panic(err)

	}
	l, err := lexer.NewLexer(absPath)
	if err != nil{
		panic(err)

	}
	grammar := GetGrammar()
	src, err := l.GetTokens()
	if err != nil{
		panic(err)
	}
	tree := ast.NewSyntaxTree(ast.NewNode(token.NewToken("", "", 0)))
	grammar["PROGRAM"].Build(&src, tree)
	return tree

}

func TestSemanticAnalyser_checkDataTypeLogicOperator(t *testing.T) {
	analyser := NewSemanticAnalyzer()
	type cases struct {
		description          string
		tree                 *ast.SyntaxTree
		expectedDataType     *token.DataType
		expectedErrorMessage string
	}
	treeOR := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.LOR, "||", 0)))
	treeOR.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "true", 0)))
	treeOR.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "false", 0)))

	treeAND := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.LAND, "&&", 0)))
	treeAND.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "true", 0)))
	treeAND.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "30", 0)))
	errorMessage := errorhandler.UnexpectedDataType(0, token.NewDataType(token.DataTypeBool, 1, 0, nil).Fmt(), token.NewDataType(token.DataTypeByte, 1, 0, nil).Fmt())

	testCases := []cases{
		{
			description:          "OR bool",
			tree:                 treeOR,
			expectedDataType:     token.NewDataType(token.DataTypeBool, 1, 0, nil),
			expectedErrorMessage: "",
		},
		{description: "AND err",
			tree:                 treeAND,
			expectedDataType:     nil,
			expectedErrorMessage: errorMessage,
		},
	}
	for _, scenario := range testCases {

		t.Run(scenario.description, func(t *testing.T) {
			dt, err := analyser.getLogicOperator(scenario.tree, "")
			assert.Equal(t, scenario.expectedDataType, dt)
			if err != nil {
				assert.Equal(t, scenario.expectedErrorMessage, err.Error())

			}

		})
	}
}
func TestSemanticAnalyser_checkDataTypeComparison(t *testing.T) {
	analyser := NewSemanticAnalyzer()
	type cases struct {
		description          string
		tree                 *ast.SyntaxTree
		expectedDataType     *token.DataType
		expectedErrorMessage string
	}
	treeLT := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.LT, token.LT, 0)))
	treeLT.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "true", 0)))
	treeLT.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "false", 0)))

	/*treeGT := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.GT, token.GT, 0)))
	treeGT.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "true", 0)))
	treeGT.Head.AddChild(treeLT.Head)
*/
	treeLTEQ := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.LTEQ, token.LTEQ, 0)))
	treeLTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "30", 0)))
	treeLTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "false", 0)))

/*	treeGTEQ := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.GTEQ, token.GTEQ, 0)))
	treeGTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "30", 0)))
	treeGTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.IDENT, "pointer", 0)))
*/
	errorMessage := errorhandler.UnexpectedDataType(0, "numeric", token.NewDataType(token.DataTypeBool, 1, 0, nil).Fmt())

	testCases := []cases{
		{
			description:          "< bool",
			tree:                 treeLT,
			expectedDataType:     token.NewDataType(token.DataTypeBool, 1, 0, nil),
			expectedErrorMessage: "",
		},
	/*	{
			description:          "> bool",
			tree:                 treeGT,
			expectedDataType:     token.NewDataType(token.DataTypeBool, 1, 0, nil),
			expectedErrorMessage: "",
		},
	*/	{description: "<= err",
			tree:                 treeLTEQ,
			expectedDataType:     nil,
			expectedErrorMessage: errorMessage,
		},
	}
	for _, scenario := range testCases {

		t.Run(scenario.description, func(t *testing.T) {
			dt, err := analyser.getComparison(scenario.tree, "")
			assert.Equal(t, scenario.expectedDataType, dt)
			if err != nil {
				assert.Equal(t, scenario.expectedErrorMessage, err.Error())

			}

		})
	}
}

func TestSemanticAnalyser_checkDataTypeEqualEqual(t *testing.T) {
	analyser := NewSemanticAnalyzer()
	type cases struct {
		description          string
		tree                 *ast.SyntaxTree
		expectedDataType     *token.DataType
		expectedErrorMessage string
	}
	treeEQEQ := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.EQEQ, token.EQEQ, 0)))
	treeEQEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "20", 0)))
	treeEQEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "30", 0)))

	treeNOTEQ := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.NOTEQ, token.NOTEQ, 0)))
	treeNOTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "true", 0)))
	treeNOTEQ.Head.AddChild(ast.NewNode(token.NewToken(token.BYTE, "20", 0)))

	treeEQEQ2 := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.EQEQ, token.EQEQ, 0)))
	treeEQEQ2.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "false", 0)))
	treeEQEQ2.Head.AddChild(treeEQEQ.Head)

	treeEQEQ3 := ast.NewSyntaxTree(ast.NewNode(token.NewToken(token.EQEQ, token.EQEQ, 0)))
	treeEQEQ3.Head.AddChild(ast.NewNode(token.NewToken(token.BOOL, "false", 0)))
	treeEQEQ3.Head.AddChild(treeNOTEQ.Head)

	errorMessage := errorhandler.DataTypesDontMatch(0, token.NewDataType(token.DataTypeBool, 1,0, nil).Fmt(),
		token.NOTEQ, token.NewDataType(token.DataTypeByte, 1,0, nil).Fmt())
	testCases := []cases{
		{
			description:          "EQEQ bool",
			tree:                 treeEQEQ,
			expectedDataType:     token.NewDataType(token.DataTypeBool, 1, 0, nil),
			expectedErrorMessage: "",
		},
			{description: "NOTEQ err",
			tree:                 treeNOTEQ,
			expectedDataType:     nil,
			expectedErrorMessage: errorMessage,
		},
		{
			description:          "EQEQ2 bool",
			tree:                 treeEQEQ2,
			expectedDataType:     token.NewDataType(token.DataTypeBool, 1, 0, nil),
			expectedErrorMessage: "",
		},
		{description: "EQEQ3 err",
			tree:                 treeEQEQ3,
			expectedDataType:     nil,
			expectedErrorMessage: errorMessage,
		},
	}
	for _, scenario := range testCases {

		t.Run(scenario.description, func(t *testing.T) {
			dt, err := analyser.getEQEQ(scenario.tree, "")
			assert.Equal(t, scenario.expectedDataType, dt)
			if err != nil {
				assert.Equal(t, scenario.expectedErrorMessage, err.Error())

			}

		})
	}

}