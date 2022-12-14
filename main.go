package main

import "github.com/NoetherianRing/c8-compiler/app"

func main (){

	compiler, err := app.NewApp()
	if err != nil{
		panic(err)
	}else{
		compiler.Program()
	}
}
