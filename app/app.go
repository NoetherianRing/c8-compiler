package app

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/lexer"
	"github.com/NoetherianRing/c8-compiler/parser"
	"github.com/NoetherianRing/c8-compiler/token"
	"os"
	"path/filepath"
	"strconv"
)

type App struct{
	lexer   *lexer.Lexer
	program *parser.NonTerminal
}

func NewApp() (*App, error){
	absPath, err := filepath.Abs(os.Args[1])
	if err != nil{
		panic (err)
	}
	l, err := lexer.NewLexer(absPath)
	if err != nil{
		return nil, err
	}
	grammar := parser.GetGrammar()


	return &App{lexer: l, program: grammar[parser.PROGRAM]}, err
}

func (app *App) Program(){
	src, err := app.lexer.GetTokens()
	if err != nil{
		panic(err)
	}
	tree := parser.NewSyntaxTree(parser.NewNode(token.NewToken("", "", 0)))
	valid := app.program.Build(&src, tree)

	if !valid{
		errorString := "syntactic error\nin line: "+ strconv.Itoa(src[0].Line) + "\nin symbol: "+ src[0].Literal
		err2:= errors.New(errorString)
		panic(err2)
	}

}
