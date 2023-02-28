package emitter

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
)

type Emitter struct{
	currentAddress uint16
	offsetStack    int //we use this field to know the last address in which we save a variable in the stack
	//globalVariables		map[string]uint16
	globalVariables map[string]uint16 //we save the address in which the global variable is stored
	variables 		map[string]int //we save the position (counting from the first variable of the list of stacks)
	// where we stored the variable. With this we can calculate the address where a variable is stored, by adding this position to
	// the address in which the stack section starts
	scope           *symboltable.Scope
	ctxNode         *ast.Node
	machineCode     [MEMORY]byte
	translateStatement     map[token.Type]func(*FunctionCtx)error
	translateOperation	   map[token.Type]func(function *FunctionCtx) (interface{}, error)
	functions       	   map[string]uint16

}

func NewEmitter(tree *ast.SyntaxTree, scope *symboltable.Scope)*Emitter{
	emitter := new(Emitter)

	emitter.globalVariables = make(map[string]uint16)
	emitter.scope = scope
	emitter.ctxNode = tree.Head
	emitter.translateStatement = make(map[token.Type]func(*FunctionCtx)error)
//	emitter.translateStatement[token.IF] = emitter._if
//	emitter.translateStatement[token.ELSE] = emitter._else
//	emitter.translateStatement[token.WHILE] = emitter._while
	emitter.translateStatement[token.EQ] = emitter.assign

	emitter.translateOperation = make(map[token.Type]func(*FunctionCtx)(interface{},error))

	emitter.currentAddress = AddressGlobalSection
	return emitter
}


//Start
func (emitter *Emitter) Start() ([MEMORY]byte, error){
	emitter.ctxNode = emitter.ctxNode.Children[0] //The tree start with a EOF node, so we move to the next one

	//we save all the global globalVariables into memory
	block := emitter.ctxNode
	for _, child := range block.Children{
		if child.Value.Type == token.LET{
			emitter.ctxNode = child
			err := emitter.globalVariableDeclaration()
			if err != nil{
				return emitter.machineCode, err
			}
		}
	}
	emitter.ctxNode = block

	//we save into memory the primitive functions
	err := emitter.primitiveFunctionsDeclaration()
	if err != nil{
		return emitter.machineCode, err
	}
	mainScope := emitter.scope
	//we save into memory all the functions (including main)
	for i, child := range block.Children{
		if child.Value.Type == token.FUNCTION{
			emitter.scope = mainScope.SubScopes[i]
			emitter.ctxNode = child
			err = emitter.functionDeclaration()
			if err != nil{
				return emitter.machineCode, err
			}

		}
	}
	emitter.scope = mainScope
	emitter.ctxNode = block


	//The stack section will start in the last available address, which is saved in the v4 and v5 registers
	v4 := byte(emitter.currentAddress & 0xFF00 >> 8)
	v5 := byte(emitter.currentAddress & 0x00FF)
	x := byte(4)
	y := byte(5)
	saveV4 := I6XKK(x, v4)
	emitter.machineCode[0] = saveV4[0]
	emitter.machineCode[1] = saveV4[1]

	saveV5 := I6XKK(y, v5)
	emitter.machineCode[2] = saveV5[0]
	emitter.machineCode[3] = saveV5[1]

	//The program will start in the main function, so we jump there
	mainAddress, ok := emitter.functions[token.MAIN]
	if !ok{
		return emitter.machineCode, errors.New(errorhandler.UnexpectedCompilerError())
	}
	jumpToStart := I1NNN(mainAddress)
	emitter.machineCode[4] = jumpToStart[0]
	emitter.machineCode[5] = jumpToStart[1]


	return emitter.machineCode, nil
}

//function declaration save the instructions of all primitive function in memory
func (emitter *Emitter) primitiveFunctionsDeclaration()error{
	return nil

}


//function declaration save all the instructions of a function in memory
func (emitter *Emitter) functionDeclaration()error{
	const ARG = 0
	const IDENT = 1
	const BLOCK = 4
	//we backup the offsetStack so we can update it after compiling the function
	offsetBackup := emitter.offsetStack

	//we store the address in which the function is saved in the map of functions
	functionName := emitter.ctxNode.Children[IDENT].Value.Literal
	emitter.functions[functionName] = emitter.currentAddress //the function starts at the current address

	mainScope := emitter.scope
	emitter.scope = emitter.scope.SubScopes[len(emitter.functions)]
	ctxAddresses := NewAddresses()

	//we save the arguments in the stack
	fn  := emitter.ctxNode
	for _, child := range fn.Children[ARG].Children{
		emitter.ctxNode = child
		err := emitter.let(ctxAddresses)
		if err != nil{
			return err
		}
	}
	emitter.ctxNode = fn

	//we declare all variables in the stack
	emitter.ctxNode = emitter.ctxNode.Children[BLOCK]
	err := emitter.declareInStack(ctxAddresses)
	emitter.ctxNode = fn
	if err != nil{
		return err
	}

	registers := NewRegisterOptimizer().optimizeRegisters(emitter.ctxNode, ctxAddresses)

	ctxFunction := NewCtxFunction(registers, ctxAddresses)

	//we write the rest of the statements in memory
	for _, child := range fn.Children[BLOCK].Children{
		emitter.ctxNode = child
		//we jump let stmts because we already save them

		translateStmt, ok :=emitter.translateStatement[emitter.ctxNode.Value.Type]
		if ok{
			err := translateStmt(ctxFunction)
			if err != nil{
				return err
			}
		}

	}

	emitter.scope = mainScope
	emitter.offsetStack = offsetBackup
	return nil
}

//globalVariableDeclaration assigns an address to a global variable and updates the current address
func (emitter *Emitter) globalVariableDeclaration() error{
	let := emitter.ctxNode
	ident := let.Children[0].Value.Literal
	symbol, ok := emitter.scope.Symbols[ident]
	if !ok{
		return errors.New(errorhandler.UnexpectedCompilerError())
	}

	switch symbol.DataType.(type) {
	case symboltable.Simple:
		emitter.globalVariables[ident] = emitter.currentAddress
		emitter.machineCode[emitter.currentAddress] = 0
		return emitter.moveCurrentAddress()
	case symboltable.Pointer:
		emitter.globalVariables[ident] = emitter.currentAddress
		emitter.machineCode[emitter.currentAddress] = 0
		err := emitter.moveCurrentAddress()
		if err != nil{
			return err
		}
		emitter.machineCode[emitter.currentAddress] = 0
		return emitter.moveCurrentAddress()

	case symboltable.Array:
		emitter.globalVariables[ident] = emitter.currentAddress
		for i := 0; i< symbol.DataType.(symboltable.Array).Length  * symbol.DataType.(symboltable.Array).SizeOfElements(); i++{
			emitter.machineCode[emitter.currentAddress] = 0
			err := emitter.moveCurrentAddress()
			if err != nil{
				return err
			}
		}
	default:
		return errors.New(errorhandler.UnexpectedCompilerError())
	}
	return nil


}

//moveCurrentAddress moves the current address by one, and if it's out of bounds of the memory it return a error
func (emitter *Emitter) moveCurrentAddress() error{
	emitter.currentAddress++
	if emitter.currentAddress > MEMORY{
		return errors.New(errorhandler.NotEnoughMemory())
	}else {
		return nil
	}

}

//declareInStack saves all variables of a function in its stack
func (emitter *Emitter) declareInStack(ctxAddresses *Addresses) error {
	backupCtxNode := emitter.ctxNode
	backupScope := emitter.scope

	iSubScope := 0

	for _, child := range emitter.ctxNode.Children{
		switch child.Value.Type {
		case token.WHILE:
				iSubScope++
				emitter.ctxNode = child.Children[1]
				emitter.scope = emitter.scope.SubScopes[iSubScope]
				ctxAddresses.AddSubAddresses()

				err := emitter.declareInStack(ctxAddresses.SubAddresses[iSubScope])
				emitter.ctxNode = backupCtxNode
				emitter.scope = backupScope
				if err != nil{
						return err
				}
		case token.IF:
			emitter.ctxNode = child.Children[1]
			emitter.scope = emitter.scope.SubScopes[iSubScope]
			ctxAddresses.AddSubAddresses()
			err := emitter.declareInStack(ctxAddresses.SubAddresses[iSubScope])
			emitter.ctxNode = backupCtxNode
			emitter.scope = backupScope
			if err != nil{
				return err
			}
		case token.ELSE:
			for j:=0; j<2;j++{
				iSubScope = iSubScope + j
				emitter.ctxNode = child.Children[j]
				emitter.scope = emitter.scope.SubScopes[iSubScope]
				ctxAddresses.AddSubAddresses()

				err := emitter.declareInStack(ctxAddresses.SubAddresses[iSubScope])
				emitter.ctxNode = backupCtxNode
				emitter.scope = backupScope
				if err != nil{
					return err
				}
			}
		case token.LET:
			err := emitter.let(ctxAddresses)
			if err != nil{
				return nil
			}

		}
	}

	emitter.ctxNode = backupCtxNode
	return nil

}

//let saves a specific variable in the stack of a function and update ctxAddresses
func (emitter *Emitter) let(ctxAddresses *Addresses) error{
	IDENT := 0
	ident := emitter.ctxNode.Children[IDENT].Value.Literal
	ctxAddresses.AddAddress(ident, emitter.offsetStack)
	symbol, ok := emitter.scope.Symbols[ident]
	if !ok{
		return errors.New(errorhandler.UnexpectedCompilerError())
	}
	size := 0
	switch symbol.DataType.(type){
	case symboltable.Simple:
		size = symbol.DataType.(symboltable.Simple).Size
	case symboltable.Pointer:
		size = symbol.DataType.(symboltable.Pointer).Size
	case symboltable.Array:
		size = symbol.DataType.(symboltable.Array).SizeOfElements() * symbol.DataType.(symboltable.Array).Length
	default:
		return errors.New(errorhandler.UnexpectedCompilerError())

	}
	if size > 16{
		for i:=size-1; i>=0; i=i-16{
			instructionsToSave := emitter.saveInStack(byte(size))
			for _, toSave := range instructionsToSave{
				emitter.machineCode[emitter.currentAddress] = toSave
				emitter.currentAddress++
			}
		}
	}else{
		instructionsToSave := emitter.saveInStack(byte(size))
		for _, toSave := range instructionsToSave{
			emitter.machineCode[emitter.currentAddress] = toSave
			emitter.currentAddress++
		}
	}

	return nil
}

//saveInStack return the instructions needed to save a variable of size "size" in the stack
func (emitter *Emitter)saveInStack(size byte)[]byte{
	instructions := make([]byte,8)

	iaxy0 := IAXY0(RegisterStackAddres1, RegisterStackAddres2) //I = (Vi << 8 | Vj)
	i6xkk := I6XKK(0, byte(emitter.offsetStack)) //v0 = offset
	ifx1e := IFX1E(0) // I = I + V0
	ifx55 := IFX55(size) // fx55 stores registers V0 through Vsize in memory starting at location I

	instructions[0] = iaxy0[0]
	instructions[1] = iaxy0[1]
	instructions[2] = i6xkk[0]
	instructions[3] = i6xkk[1]
	instructions[4] = ifx1e[0]
	instructions[5] = ifx1e[1]
	instructions[6] = ifx55[0]
	instructions[7] = ifx55[1]

	return instructions

}

//assign translate the assign statement to opcodes and write it in emitter.machineCode
func (emitter *Emitter)assign(functionCtx *FunctionCtx)error{
	const SAVEIN = 0
	const TOSAVE = 1

	assignBackup := emitter.ctxNode
	emitter.ctxNode = emitter.ctxNode.Children[TOSAVE]
	//depending on the data type of the right side of the assign translateOperation save it in v0 or v0 and v1
	datatype, err :=emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err!= nil{
		return err
	}
	emitter.ctxNode = assignBackup

	if emitter.ctxNode.Children[SAVEIN].Value.Type == token.IDENT{
		ident := emitter.ctxNode.Children[SAVEIN].Value.Literal
		reference, referenceExists := functionCtx.Addresses.GetReference(ident)

		if !referenceExists {
			_, isAGlobalVariable := emitter.globalVariables[ident]
			if !isAGlobalVariable{
				return errors.New(errorhandler.UnexpectedCompilerError())
			}else{
				//if the reference is to a global variable we update I = address
				err := emitter.saveGlobalReferenceAddressInI(2, 3)
				if err != nil{
					return err
				}
				//and then we store registers V0 through Vsize in memory starting at location I

				ifx55 := IFX55(byte(symboltable.GetSize(datatype)))

				emitter.machineCode[emitter.currentAddress] = ifx55[0]
				err =emitter.moveCurrentAddress()
				if err != nil{
					return err
				}

				emitter.machineCode[emitter.currentAddress] = ifx55[1]
				err =emitter.moveCurrentAddress()
				if err != nil{
					return err
				}
				return nil

			}


		}

		//if the left side of the assign is in a register x we just write vx = v0 (8X00) (because to be in a register it must be a simple)
		indexRegister, referenceIsInRegister := functionCtx.Registers.guide[reference]
		if referenceIsInRegister {
			i8xy0 := I8XY0(byte(indexRegister), 0)
			emitter.machineCode[emitter.currentAddress] = i8xy0[0]
			err := emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
			emitter.machineCode[emitter.currentAddress] = i8xy0[1]
			err = emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
			return nil
		}
	}

	//if not, we search the address of the left side of the assign and we save it in I

	if emitter.ctxNode.Value.Type == token.IDENT{
		err := emitter.saveStackReferenceAddressInI(2, functionCtx)
		if err != nil{
			return err
		}
	}else{
		//err := emitter.saveDereferenceAddressInI(2, 3, functionCtx)
		//because to save a dereference address in i we need v0 and v1 we need to backup v0 (and v1) in v2 (and v3)
		i8xy0 := I8XY0(2, 0)
		emitter.machineCode[emitter.currentAddress] = i8xy0[0]
		err := emitter.moveCurrentAddress()
		if err!=nil{
			return err
		}
		emitter.machineCode[emitter.currentAddress] = i8xy0[1]
		err = emitter.moveCurrentAddress()
		if err!=nil{
			return err
		}
		if symboltable.GetSize(datatype) > 1{
			i8xy0 := I8XY0(3, 1)
			emitter.machineCode[emitter.currentAddress] = i8xy0[0]
			err := emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
			emitter.machineCode[emitter.currentAddress] = i8xy0[1]
			err = emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
		}
		err = emitter.saveDereferenceAddressInI(functionCtx)
		if err != nil{
			return err
		}
		//because to save a vx (and vy) in memory we need them in v0 (and v1), we save them there again
		i8xy0 = I8XY0(0, 2)
		emitter.machineCode[emitter.currentAddress] = i8xy0[0]
		err = emitter.moveCurrentAddress()
		if err!=nil{
			return err
		}
		emitter.machineCode[emitter.currentAddress] = i8xy0[1]
		err = emitter.moveCurrentAddress()
		if err!=nil{
			return err
		}
		if symboltable.GetSize(datatype) > 1{
			i8xy0 := I8XY0(1, 3)
			emitter.machineCode[emitter.currentAddress] = i8xy0[0]
			err := emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
			emitter.machineCode[emitter.currentAddress] = i8xy0[1]
			err = emitter.moveCurrentAddress()
			if err!=nil{
				return err
			}
		}
		emitter.ctxNode = assignBackup.Children[SAVEIN]
	}


	//now we store registers V0 through Vsize in memory starting at location I.
	ifx55 := IFX55(byte(symboltable.GetSize(datatype)))

	emitter.machineCode[emitter.currentAddress] = ifx55[0]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}

	emitter.machineCode[emitter.currentAddress] = ifx55[1]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}

	return nil
}

//saveGlobalReferenceAddressInI saves the address of a global variable in I using the register x and y
func (emitter *Emitter)saveGlobalReferenceAddressInI(x byte, y byte) error{
	ident := emitter.ctxNode.Value.Literal
	address := emitter.globalVariables[ident]
	i6xkk1 := I6XKK(x, byte(address<<8))
	i6xkk2 := I6XKK(y, byte(address))
	iaxy0 := IAXY0(x,y)
	emitter.machineCode[emitter.currentAddress] = i6xkk1[0]
	err := emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.machineCode[emitter.currentAddress] = i6xkk1[1]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.machineCode[emitter.currentAddress] = i6xkk2[0]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.machineCode[emitter.currentAddress] = i6xkk2[1]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}

	emitter.machineCode[emitter.currentAddress] = iaxy0[0]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.machineCode[emitter.currentAddress] = iaxy0[1]
	err =emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	return nil

}

//saveDereferenceAddressInI save the address of a dereference in I using the registers 0 and 1
func (emitter *Emitter)saveDereferenceAddressInI(functionCtx *FunctionCtx) error {
	return nil
}

//saveStackReferenceAddressInI save the address of a reference saved in the stack in I using the register x
func (emitter *Emitter)saveStackReferenceAddressInI( x byte, functionCtx *FunctionCtx)error{
	ident := emitter.ctxNode.Value.Literal
	reference, ok := functionCtx.Addresses.References[ident]
	if !ok{
		return errors.New(errorhandler.UnexpectedCompilerError())
	}

	//we set I = address position 0 of stack
	iaxy0 := IAXY0(RegisterStackAddres1, RegisterStackAddres2)

	emitter.machineCode[emitter.currentAddress] = iaxy0[0]
	err := emitter.moveCurrentAddress()
	if err != nil{
		return nil
	}
	emitter.machineCode[emitter.currentAddress] = iaxy0[1]
	err = emitter.moveCurrentAddress()
	if err != nil{
		return nil
	}
	if reference.positionStack <= 255{
		i6xkk := I6XKK(x, byte(reference.positionStack))
		ifx1e := IFX1E(x)

		emitter.machineCode[emitter.currentAddress] = i6xkk[0]
		err := emitter.moveCurrentAddress()
		if err != nil{
			return err
		}
		emitter.machineCode[emitter.currentAddress] = i6xkk[1]
		err = emitter.moveCurrentAddress()
		if err != nil{
			return err
		}

		emitter.machineCode[emitter.currentAddress] = ifx1e[0]
		err = emitter.moveCurrentAddress()
		if err != nil{
			return err
		}

		emitter.machineCode[emitter.currentAddress] = ifx1e[1]
		err = emitter.moveCurrentAddress()
		if err != nil{
			return err
		}

	}else{
		position := reference.positionStack

		for  position>255{
			i6xkk := I6XKK(x, 255)
			ifx1e := IFX1E(x)

			emitter.machineCode[emitter.currentAddress] = i6xkk[0]
			err := emitter.moveCurrentAddress()
			if err != nil{
				return err
			}
			emitter.machineCode[emitter.currentAddress] = i6xkk[1]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}

			emitter.machineCode[emitter.currentAddress] = ifx1e[0]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}

			emitter.machineCode[emitter.currentAddress] = ifx1e[1]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}
			position = position - 255
		}
		if position > 0{
			i6xkk := I6XKK(x, byte(position))
			ifx1e := IFX1E(x)
			emitter.machineCode[emitter.currentAddress] = i6xkk[0]
			err := emitter.moveCurrentAddress()
			if err != nil{
				return err
			}
			emitter.machineCode[emitter.currentAddress] = i6xkk[1]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}

			emitter.machineCode[emitter.currentAddress] = ifx1e[0]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}

			emitter.machineCode[emitter.currentAddress] = ifx1e[1]
			err = emitter.moveCurrentAddress()
			if err != nil{
				return err
			}
		}

	}
	return nil

}