package app

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	emitter2 "github.com/NoetherianRing/c8-compiler/emitter"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/semanticAnalyzer"
	"github.com/NoetherianRing/c8-compiler/syntacticanalyzer"
	"github.com/NoetherianRing/c8-compiler/token"
	"os"
	"path/filepath"
)

type App struct{
	sourceFilePath string
	romFilePath string
	program *syntacticanalyzer.NonTerminal
}

func NewApp(sourceFilePath string, romFilePath string) (*App, error){
	var err error
	app := new(App)
	app.sourceFilePath, err = filepath.Abs(sourceFilePath)
	if err != nil{
		return nil, err
	}
	app.romFilePath, err = filepath.Abs(romFilePath)
	if err != nil{
		return nil, err
	}
	grammar := syntacticanalyzer.GetGrammar()
	app.program = grammar[syntacticanalyzer.PROGRAM]


	return app, err
}

func (app *App) Program(){

	l, err := lexer.NewLexer(app.sourceFilePath)
	if err != nil{
		panic(err)
	}

	src, err :=l.GetTokens()
	if err != nil{
		panic(err)
	}

	tree := ast.NewSyntaxTree(ast.NewNode(token.NewToken("", "", 0)))
	valid := app.program.Build(&src, tree)
	if !valid{
		err = errors.New(errorhandler.SyntaxError())
		panic(err)
	}

	semantic := semanticAnalyzer.NewSemanticAnalyzer(tree)
	scope, err := semantic.Start()
	if err != nil{
		panic(err)
	}
	emitter := emitter2.NewEmitter(tree, scope)
	machineCode, err := emitter.Start()
	if err != nil{
		panic(err)
	}
	f, err := os.Create(app.romFilePath)
	if err !=nil{
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(machineCode)

	if err != nil{
		panic(err)

	}

}
