package app

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/syntacticanalyzer"
	"github.com/NoetherianRing/c8-compiler/token"
	"path/filepath"
	"strconv"
)

type App struct{
	lexer   *lexer.Lexer
	program *syntacticanalyzer.NonTerminal
}

func NewApp(path string) (*App, error){
	absPath, err := filepath.Abs(path)
	if err != nil{
		panic (err)
	}
	l, err := lexer.NewLexer(absPath)
	if err != nil{
		return nil, err
	}
	grammar := syntacticanalyzer.GetGrammar()


	return &App{lexer: l, program: grammar[syntacticanalyzer.PROGRAM]}, err
}

func (app *App) Program(){
	src, err := app.lexer.GetTokens()
	if err != nil{
		panic(err)
	}
	tree := ast.NewSyntaxTree(ast.NewNode(token.NewToken("", "", 0)))
	valid := app.program.Build(&src, tree)

	if !valid{
		errorString := "syntactic errorhandler\nin line: "+ strconv.Itoa(src[0].Line) + "\nin symbol: "+ src[0].Literal
		err2:= errors.New(errorString)
		panic(err2)
	}

}
