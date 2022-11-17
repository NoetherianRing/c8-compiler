package app

import (
	"github.com/NoetherianRing/c8-compiler/lexer"
	"os"
)

type App struct{
	l *lexer.Lexer
}

func NewApp() (*App, error){
	l, err := lexer.NewLexer(os.Args[1])
	if err != nil{
		return nil, err
	}
	return &App{l: l}, err
}

func (app *App) Program(){
	/*tokens, err := app.l.GetTokens()
	if err != nil{
		panic(err) //TODO: Esto lo tengo que mostrar al final creo, tipo en el emmiter
	}
*/
}
