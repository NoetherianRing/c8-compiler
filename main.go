package main

import (
	"github.com/NoetherianRing/c8-compiler/app"
	"os"
)

func main (){

	compiler, err := app.NewApp(os.Args[1], os.Args[2])
	if err != nil{
		panic(err)
	}
	compiler.Program()

}
