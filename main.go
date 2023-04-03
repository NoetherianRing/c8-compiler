package main

import (
	"github.com/NoetherianRing/c8-compiler/app"
	"os"
)

func main() {
	const inputPathArg = 1
	const outputPathArg = 2
	compiler, err := app.NewApp(os.Args[inputPathArg], os.Args[outputPathArg])
	if err != nil {
		panic(err)
	}
	compiler.Program()

}
