package emitter

import (
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/token"
	"sort"
)

type RegisterOptimizer struct{
	registers *RegistersGuide
	count 	map[*Reference]int //it is used to count how many times a reference is used in a function

}

type RegistersGuide struct{
	guide 	map[*Reference]int //it tells what addresses are saved in registers, and in which registers they are saved

}

func  NewRegisterOptimizer()*RegisterOptimizer{

	optimizer := new(RegisterOptimizer)
	optimizer.registers = newRegisters()
	optimizer.count = make(map[*Reference]int)
	return optimizer
}
func  newRegisters()*RegistersGuide {
	registers := new(RegistersGuide)
	registers.guide =  make(map[*Reference]int)
	return registers
}

//optimizeRegisters receive a node and an StackReferences, counts how many times each variable appears in the given context,
//decide which ones would be optimus to save in registers, and returns a RegistersGuide which contains a map whith the reference of the variable
//as key, and the register in which is saved as value
func (optimizer *RegisterOptimizer)optimizeRegisters(ctxNode *ast.Node,ctxAddresses *StackReferences) *RegistersGuide {

	optimizer.toCount(ctxNode, ctxAddresses)
	if len(optimizer.count) != 0{
		keys := make([]*Reference, 0, len(optimizer.count))


		for key := range optimizer.count {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(i, j int) bool{
			return optimizer.count[keys[i]] < optimizer.count[keys[j]]
		})
		for i := 4; i <=13; i++{ //the first 4 registers and the last 3 are used as aux in operations, so we don't save variables there
			if i - 4 > len(keys){
				return optimizer.registers
			}
			optimizer.registers.guide[keys[i-4]] = i
		}


	}
	return optimizer.registers

}

func (optimizer *RegisterOptimizer)toCount(ctxNode *ast.Node,ctxAddresses *StackReferences) {

	if ctxNode.Value.Type == token.IDENT{
		reference, ok := ctxAddresses.References[ctxNode.Value.Literal]
		if ok{
			optimizer.count[reference] +=1
		}
	}else{
		subScopesFounded := 0
		for _, child := range ctxNode.Children{
			if child.Value.Type == token.RBRACE{
				nextCtxAddresses := ctxAddresses.SubReferences[subScopesFounded]
				subScopesFounded++
				optimizer.toCount(child, nextCtxAddresses)
			}else{
				optimizer.toCount(child, ctxAddresses)
			}
		}
	}
}