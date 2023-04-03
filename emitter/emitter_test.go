package emitter

import (
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/semanticAnalyzer"
	"github.com/NoetherianRing/c8-compiler/syntacticanalyzer"
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestStart(t *testing.T) {
	type cases struct {
		description string
		testPathTxt string
		testPathRom string
		err         error
	}
	const numberOfValidTests = 41
	testCases := make([]cases, 0)
	for i := 0; i < numberOfValidTests; i++ {
		pathTxt := "../fixtures/emitter/c8-lang/test" + strconv.Itoa(i+1) + ".txt"
		absPathTxt, err := filepath.Abs(pathTxt)
		pathRom := "../fixtures/emitter/roms/test" + strconv.Itoa(i+1) + ".ch8"
		absPathRom, err := filepath.Abs(pathRom)

		if err != nil {
			assert.Error(t, err)
		}
		testCases = append(testCases, cases{
			description: "test" + strconv.Itoa(i+1),
			testPathTxt: absPathTxt,
			testPathRom: absPathRom,
			err:         nil,
		})

	}
	grammar := syntacticanalyzer.GetGrammar()
	program := grammar[syntacticanalyzer.PROGRAM]
	for _, scenario := range testCases {
		l, err := lexer.NewLexer(scenario.testPathTxt)
		assert.NoError(t, err)

		tokens, err := l.GetTokens()
		assert.NoError(t, err)

		tree := ast.NewSyntaxTree(ast.NewNode(token.NewToken("", "", 0)))
		valid := program.Build(&tokens, tree)
		assert.True(t, valid, "invalid syntax")
		semantic := semanticAnalyzer.NewSemanticAnalyzer(tree)
		scope, err := semantic.Start()
		assert.NoError(t, err)
		emitter := NewEmitter(tree, scope)
		machineCode, err := emitter.Start()
		assert.NoError(t, err)
		rom, err := os.ReadFile(scenario.testPathRom)
		assert.NoError(t, err)
		assert.NotEqual(t, t, machineCode, rom)
	}

}
