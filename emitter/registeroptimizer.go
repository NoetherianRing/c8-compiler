package emitter

import "github.com/NoetherianRing/c8-compiler/ast"

type RegisterOptimizer struct{

}

type Registers struct{
	registers [16]byte          //it saves the value of the registers
	guide 	map[*Reference]int //it tells what addresses are saved in registers, and in which registers they are saved
}

func  NewRegisterOptimizer()*RegisterOptimizer{
	return new(RegisterOptimizer)
}
//TODO
func (optimizer *RegisterOptimizer)optimizeRegisters(ctxNode *ast.Node,ctxAddresses *Addresses) *Registers{
	return nil
}