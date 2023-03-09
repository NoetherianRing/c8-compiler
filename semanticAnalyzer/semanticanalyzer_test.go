package semanticAnalyzer

import (
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/syntacticanalyzer"
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strconv"
	"testing"
)

func TestStart(t *testing.T) {
	type cases struct{
		description 	string
		testPath 		string
		err  			error
	}
	const numberOfValidTests = 4
	testCases := make([]cases,0)
	for i:=0; i<numberOfValidTests; i++{
		path := "../fixtures/semantic/valid/valid_test" + strconv.Itoa(i) +".text"
		absPath, err := filepath.Abs(path)
		if err != nil{
			assert.Error(t, err)
		}
		testCases = append(testCases, cases{
			description: "valid_test" + strconv.Itoa(i),
			testPath: absPath,
			err: nil,
		})


	}
	grammar := syntacticanalyzer.GetGrammar()
	program := grammar[syntacticanalyzer.PROGRAM]
	for _, scenario := range testCases{
		l, err := lexer.NewLexer(scenario.testPath)

		assert.NoError(t, err)

		tokens, err:=l.GetTokens()
		assert.NoError(t, err)

		tree := ast.NewSyntaxTree(ast.NewNode(token.NewToken("", "", 0)))
		valid := program.Build(&tokens, tree)
		assert.True(t, valid, "invalid syntax")
		semantic :=NewSemanticAnalyzer(tree)
		_, err = semantic.Start()
		assert.NoError(t, err)

	}

}