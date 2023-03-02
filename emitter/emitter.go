package emitter

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
	"strconv"
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
	machineCode    		   [MEMORY]byte
	translateStatement     map[token.Type]func(*FunctionCtx)error
	translateOperation	   map[token.Type]func(function *FunctionCtx) (int, error)
	functions       	   map[string]uint16
	lastIndexSubScope	   int //in the context of a scope, tells the numbers of sub-scopes already written in machine code

}

func NewEmitter(tree *ast.SyntaxTree, scope *symboltable.Scope)*Emitter{
	emitter := new(Emitter)

	emitter.globalVariables = make(map[string]uint16)
	emitter.scope = scope
	emitter.lastIndexSubScope = 0
	emitter.ctxNode = tree.Head
	emitter.translateStatement = make(map[token.Type]func(*FunctionCtx)error)
//	emitter.translateStatement[token.IF] = emitter._if
//	emitter.translateStatement[token.ELSE] = emitter._else
//	emitter.translateStatement[token.WHILE] = emitter._while
	emitter.translateStatement[token.EQ] = emitter.assign

	emitter.translateOperation = make(map[token.Type]func(*FunctionCtx)(int,error))
	emitter.translateOperation[token.DOLLAR] = emitter.address

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
	emitter.lastIndexSubScope = 0

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
	emitter.offsetStack += int(size)
	return instructions

}

//assign translates the assign statement to opcodes and write it in emitter.machineCode
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
				_, err := emitter.saveGlobalReferenceAddressInI(2, 3)
				if err != nil{
					return err
				}
				//and then we store registers V0 through Vsize in memory starting at location I

				ifx55 := IFX55(byte(symboltable.GetSize(datatype)))

				err =emitter.saveOpcode(ifx55)

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
			err =emitter.saveOpcode(i8xy0)

			if err!=nil{
				return err
			}
			return nil
		}
	}

	//if not, we search the address of the left side of the assign and we save it in I

	if emitter.ctxNode.Value.Type == token.IDENT{
		_,err := emitter.saveStackReferenceAddressInI(2, functionCtx)
		if err != nil{
			return err
		}
	}else{
		//err := emitter.saveDereferenceAddressInI(2, 3, functionCtx)
		//because to save a dereference address in i we need v0 and v1 we need to backup v0 (and v1) in v2 (and v3)
		i8xy0 := I8XY0(2, 0)
		err =emitter.saveOpcode(i8xy0)
		if err!=nil{
			return err
		}
		if symboltable.GetSize(datatype) > 1{
			i8xy0 := I8XY0(3, 1)
			err =emitter.saveOpcode(i8xy0)

			if err!=nil{
				return err
			}
		}
		_, err = emitter.saveDereferenceAddressInI(functionCtx)
		if err != nil{
			return err
		}
		//because to save a vx (and vy) in memory we need them in v0 (and v1), we save them there again
		i8xy0 = I8XY0(0, 2)
		err =emitter.saveOpcode(i8xy0)

		if err!=nil{
			return err
		}
		if symboltable.GetSize(datatype) > 1{
			i8xy0 := I8XY0(1, 3)
			err =emitter.saveOpcode(i8xy0)

			if err!=nil{
				return err
			}
		}
		emitter.ctxNode = assignBackup.Children[SAVEIN]
	}


	//now we store registers V0 through Vsize in memory starting at location I.
	ifx55 := IFX55(byte(symboltable.GetSize(datatype)))

	err =emitter.saveOpcode(ifx55)

	if err != nil{
		return err
	}

	return nil
}

//_if translates the if statement to opcodes and write it in emitter.machineCode
func (emitter *Emitter) _if(functionCtx *FunctionCtx) error{
	const CONDITION = 0
	const BLOCK = 1

	backup := emitter.ctxNode
	//first we write in v0 the result of the condition
	emitter.ctxNode = emitter.ctxNode.Children[CONDITION]
	_, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return err
	}
	i3xkk := I3XKK(0, True) //if v0 = true we skip the next instruction
	err = emitter.saveOpcode(i3xkk)
	if err != nil{
		return err
	}
	//the next instruction is going to be a jump to the memory address after the block
	//because we don't know this address yet, we save the current address to write the opcode later
	lineAfterCondition := emitter.currentAddress
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.ctxNode = backup
	//we write all the opcodes of the block
	emitter.ctxNode = emitter.ctxNode.Children[BLOCK]
	err = emitter.block(functionCtx)
	if err != nil{
		return err
	}
	//then we write the jump after the condition
	i1nnn := I1NNN(emitter.currentAddress)

	emitter.machineCode[lineAfterCondition] = i1nnn[0]
	emitter.machineCode[lineAfterCondition+1] = i1nnn[1]

	return nil
}

//_else translates the else statement to opcodes and write it in emitter.machineCode
func (emitter *Emitter) _else(functionCtx *FunctionCtx) error{
	const CONDITION = 0
	const IFBLOCK = 1
	const ELSEBLOCK = 2

	backup := emitter.ctxNode
	//first we write in v0 the result of the condition
	emitter.ctxNode = emitter.ctxNode.Children[CONDITION]
	_, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return err
	}
	i3xkk := I3XKK(0, True) //if v0 = true we skip the next instruction
	err = emitter.saveOpcode(i3xkk)
	if err != nil{
		return err
	}
	//the next instruction is going to be a jump to the memory address after the if block
	//because we don't know this address yet, we save the current address to write the opcode later
	lineAfterCondition := emitter.currentAddress
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.ctxNode = backup
	//we write all the opcodes of the if block
	emitter.ctxNode = emitter.ctxNode.Children[IFBLOCK]
	err = emitter.block(functionCtx)
	if err != nil{
		return err
	}
	//if we execute the if block we want to jump the else block,
	//because we don't have the address to jump yet, we save the current address to write it later
	lineAfterIf := emitter.currentAddress
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.ctxNode = backup

	//then we write the jump after the condition
	i1nnn := I1NNN(emitter.currentAddress)

	emitter.machineCode[lineAfterCondition] = i1nnn[0]
	emitter.machineCode[lineAfterCondition+1] = i1nnn[1]

	//we write all the opcodes of the else block
	emitter.ctxNode = emitter.ctxNode.Children[ELSEBLOCK]
	err = emitter.block(functionCtx)
	if err != nil{
		return err
	}
	//then we write the jump after the if block
	i1nnn2 := I1NNN(emitter.currentAddress)

	emitter.machineCode[lineAfterIf] = i1nnn2[0]
	emitter.machineCode[lineAfterIf+1] = i1nnn2[1]
	return nil
}
//_while translates the else statement to opcodes and write it in emitter.machineCode
func (emitter *Emitter)_while(functionCtx *FunctionCtx) error {
	const CONDITION = 0
	const BLOCK = 1

	//we save the initial address to jump in every iteration
	initial := emitter.currentAddress
	backup := emitter.ctxNode

	emitter.ctxNode = emitter.ctxNode.Children[CONDITION]
	_, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return err
	}
	i3xkk := I3XKK(0, True) //if v0 = true we skip the next instruction
	err = emitter.saveOpcode(i3xkk)
	if err != nil{
		return err
	}
	//the next instruction is going to be a jump to the memory address after the while
	//because we don't know this address yet, we save the current address to write the opcode later
	lineAfterCondition := emitter.currentAddress
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.ctxNode = backup

	emitter.ctxNode = emitter.ctxNode.Children[BLOCK]
	err = emitter.block(functionCtx)
	if err != nil{
		return err
	}
	jumpToInitial := I1NNN(initial) //after executing the block, we jump to the address of the condition
	err =emitter.saveOpcode(jumpToInitial)
	if err != nil{
		return err
	}

	//then we write the jump after the condition
	jumpWhile := I1NNN(emitter.currentAddress)

	emitter.machineCode[lineAfterCondition] = jumpWhile[0]
	emitter.machineCode[lineAfterCondition+1] = jumpWhile[1]
	return nil
}

//block translates a block  to opcodes and write it in emitter.machineCode, it also handle the scope
func (emitter *Emitter)block(functionCtx *FunctionCtx) error {
	lastIndexSubScopeBackup := emitter.lastIndexSubScope
	scopeBackup := emitter.scope
	emitter.scope = emitter.scope.SubScopes[emitter.lastIndexSubScope]
	addressesBackup := functionCtx.Addresses
	functionCtx.Addresses = functionCtx.Addresses.SubAddresses[emitter.lastIndexSubScope]
	emitter.lastIndexSubScope = 0

	block := emitter.ctxNode

	for _, child :=range  block.Children{
		//we jump let statements because they were already executed in the declaration of a function
		if child.Value.Type != token.LET{
			emitter.ctxNode = child
			err := emitter.translateStatement[emitter.ctxNode.Value.Type](functionCtx)
			if err != nil{
				return err
			}
		}
	}

	emitter.lastIndexSubScope = lastIndexSubScopeBackup + 1
	emitter.scope = scopeBackup
	functionCtx.Addresses = addressesBackup
	return nil
}

//voidCall translates a void call to opcodes and write it in emitter.machineCode
func (emitter *Emitter)voidCall(functionCtx *FunctionCtx)error{
	_,err := emitter.call(functionCtx)
	return err
}

//call translates a call to opcodes and write it in emitter.machineCode, return the size of the datatype it returns and an error
func (emitter *Emitter)call(functionCtx *FunctionCtx)(int, error){
  	const IDENT = 0
  	const PARAMS = 1
	ident := emitter.ctxNode.Children[IDENT].Value.Literal
	//we first backup all registers of the current function in a the stack
	iaxy0 := IAXY0(RegisterStackAddres1, RegisterStackAddres2)
	backupOffset := emitter.offsetStack
	emitter.offsetStack += 16
	ifx55 := IFX55(16)

	err := emitter.saveOpcode(iaxy0)
	if err != nil{
		return 0, err
	}
	err = emitter.executeFX1ESafe(0, emitter.offsetStack)
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(ifx55)
	if err != nil{
		return 0, err
	}

	//then we obtain the value of the parameters and we save them in registers:
	if len(emitter.ctxNode.Children) > 1{

		emitter.ctxNode = emitter.ctxNode.Children[PARAMS]
		i := 0
		for emitter.ctxNode.Value.Type == token.COMMA{
			backupComma := emitter.ctxNode
			emitter.ctxNode = emitter.ctxNode.Children[0]
			//we save in v0(and maybe v1) the value of the parameter being analyzed
			size, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
			i8XY0 := I8XY0(byte(i), 0)
			i++
			err = emitter.saveOpcode(i8XY0)
			if err != nil{
				return 0, err
			}
			if size > 1{
				i8XY0 := I8XY0(byte(i), 1)
				i++
				err := emitter.saveOpcode(i8XY0)
				if err != nil{
					return 0, err
				}
			}
			emitter.ctxNode = backupComma
			emitter.ctxNode = emitter.ctxNode.Children[1]
		}
		for emitter.ctxNode.Value.Type == token.COMMA{
			emitter.ctxNode = emitter.ctxNode.Children[0]
			//we save in v0(and maybe v1) the value of the parameter being analyzed
			size, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
			i8XY0 := I8XY0(byte(i), 0)
			i++
			err = emitter.saveOpcode(i8XY0)
			if err != nil{
				return 0, err
			}
			if size > 1{
				i8XY0 := I8XY0(byte(i), 1)
				i++
				err := emitter.saveOpcode(i8XY0)
				if err != nil{
					return 0, err
				}
			}

		}

	}

	//we call the function
	fnAddress, _ := emitter.functions[ident]
	i2nnn := I2NNN(fnAddress)
	err = emitter.saveOpcode(i2nnn)
	if err != nil{
		return 0, err
	}

	//when we return, we save again the registers in memory
	emitter.offsetStack = backupOffset
	err = emitter.saveOpcode(iaxy0)
	if err != nil{
		return 0, err
	}
	err = emitter.executeFX1ESafe(0, emitter.offsetStack)
	if err != nil{
		return 0, err
	}
	ifx65 := IFX65(16)
	err = emitter.saveOpcode(ifx65)
	if err != nil{
		return 0, err
	}

	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType.(symboltable.Function).Return)
	return size, nil
}

//literal save a byte in v0. Return the size of the byte datatype (1) and an error
func (emitter *Emitter) literal(functionCtx *FunctionCtx) (int, error) {
	kk,_ := strconv.Atoi(emitter.ctxNode.Value.Literal)
	i6kk := I6XKK(0, byte(kk))
	err := emitter.saveOpcode(i6kk)
	if err != nil{
		return 0, err
	}
	return 1, nil
}

//literal save a boolean in v0. Return the size of the boolean datatype (1) and an error
func (emitter *Emitter)boolean(functionCtx *FunctionCtx) (int, error) {
	var kk byte
	if emitter.ctxNode.Value.Literal == token.TRUE{
		kk = 0xFF
	}else{
		kk = 0x00
	}
	i6kk := I6XKK(0, kk)
	err := emitter.saveOpcode(i6kk)
	if err != nil{
		return 0, err
	}
	return 1, nil
}

func (emitter *Emitter) multiply(functionCtx *FunctionCtx) (int, error) {

}

//index save in v0 (and maybe v1) the value of a dereference.
//Returns the size of the datatype of the dereference and a error
func (emitter *Emitter)index(functionCtx *FunctionCtx) (int, error){
	return emitter.saveDereferenceInRegisters(functionCtx)
}

//asterisk multiply registers or save a dereference in registers, depending on the context
//it return the size of the datatype obtained at the end of the operation and an error
func (emitter *Emitter)asterisk(functionCtx *FunctionCtx)(int, error){
	if len(emitter.ctxNode.Children) == 1{
		return emitter.saveDereferenceInRegisters(functionCtx)
	}else{
		return emitter.multiply(functionCtx)
	}

}

//saveDereferenceInRegisters save in v0 (and maybe v1) the value of a dereference.
//Returns the size of the datatype of the dereference and a error
func (emitter *Emitter) saveDereferenceInRegisters(functionCtx *FunctionCtx) (int, error) {
	size, err := emitter.saveDereferenceAddressInI(functionCtx)
	if err != nil {
		return 0, err
	}
	ifx65 := IFX65(byte(size))
	err = emitter.saveOpcode(ifx65)
	if err != nil {
		return 0, err
	}
	return size, nil
}
//ident save in v0 (and maybe v1) the value of a reference.
//Returns the size of the datatype of the reference and a error
func (emitter *Emitter)ident(functionCtx *FunctionCtx) (int, error){
	ident := emitter.ctxNode.Value.Literal
	var err error
	var size int
	_, isGlobalReference := emitter.globalVariables[ident]
	if isGlobalReference{
		size, err = emitter.saveGlobalReferenceAddressInI(0,1)
		if err != nil{
			return 0, err
		}
	}else{
		size, err = emitter.saveStackReferenceAddressInI(0,functionCtx)
		if err != nil{
			return 0, err
		}
	}
	ifx65 := IFX65(byte(size))
	err = emitter.saveOpcode(ifx65)
	if err != nil{
		return 0, err
	}
	return size, nil
}

//address save the address of its children in v0 and v1, return a error if needed, and the size of the pointer
func (emitter *Emitter)address(functionCtx *FunctionCtx)(int, error){
	emitter.ctxNode = emitter.ctxNode.Children[0]
	//first we save the address in I
	if emitter.ctxNode.Value.Type == token.IDENT{
		ident := emitter.ctxNode.Value.Literal
		_, isGlobalReference := emitter.globalVariables[ident]
		if isGlobalReference{
			_, err := emitter.saveGlobalReferenceAddressInI(0,1)
			if err != nil{
				return 0, err
			}
		}else{
			_, err := emitter.saveStackReferenceAddressInI(0,functionCtx)
			if err != nil{
				return 0, err
			}
		}
	}else{
		_,err := emitter.saveDereferenceAddressInI(functionCtx)
		if err != nil{
			return 0, err
		}
	}
	//then we save i in v0 and v1
	ibxy0 := IBXY0(0,1)
	err := emitter.saveOpcode(ibxy0)
	if err != nil{
		return 0, err
	}
	return 	2, nil
}

//saveGlobalReferenceAddressInI saves the address of a global variable in I using the register x and y
//Returns the size of the reference it points to and an error
func (emitter *Emitter)saveGlobalReferenceAddressInI(x byte, y byte) (int, error){
	ident := emitter.ctxNode.Value.Literal
	address := emitter.globalVariables[ident]
	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType)
	i6xkk1 := I6XKK(x, byte(address<<8))
	i6xkk2 := I6XKK(y, byte(address))
	iaxy0 := IAXY0(x,y)

	err :=emitter.saveOpcode(i6xkk1)
	if err != nil{
		return 0,err
	}

	err =emitter.saveOpcode(i6xkk2)
	if err != nil{
		return 0,err
	}

	err =emitter.saveOpcode(iaxy0)
	if err != nil{
		return 0,err
	}
	return size,nil

}

//saveDereferenceAddressInI save the address of a dereference in I using the registers 0 and 1.
//Returns the size of the reference it points to and an error
func (emitter *Emitter)saveDereferenceAddressInI(functionCtx *FunctionCtx) (int, error) {
	backup := emitter.ctxNode
	//we save the address of the leaf in I:
	leaf := GetLeafByRight(emitter.ctxNode)
	emitter.ctxNode = leaf
	leafIdent := emitter.ctxNode.Value.Literal
	_, isInStack := functionCtx.Addresses.References[leafIdent]
	if !isInStack{
		_, isInGlobalMemory := emitter.globalVariables[leafIdent]
		if !isInGlobalMemory{
			return 0,errors.New(errorhandler.UnexpectedCompilerError())
		}
		_, err := emitter.saveGlobalReferenceAddressInI(0,1)
		if err != nil{
			return 0, err
		}
	}else{
		_, err := emitter.saveStackReferenceAddressInI(0, functionCtx)
		if err != nil{
			return 0,err
		}
	}
	symbol, ok := emitter.scope.Symbols[leafIdent]

	if !ok{
		return 0,errors.New(errorhandler.UnexpectedCompilerError())
	}
	datatype := symbol.DataType

	emitter.ctxNode = backup

	for emitter.ctxNode!=leaf{

		switch emitter.ctxNode.Value.Type{
		//if we are analyzing a ], we add the index of the array to I to set the address of the next referenced element in I
		case token.RBRACKET:
			index, err := strconv.Atoi(emitter.ctxNode.Children[0].Value.Literal)
			if err != nil{
				return 0,errors.New(errorhandler.UnexpectedCompilerError())
			}
			err = emitter.executeFX1ESafe(0, index*symboltable.GetSize(datatype))
			datatype = datatype.(symboltable.Array).Of
			emitter.ctxNode = emitter.ctxNode.Children[1]

		//if we are analyzing a *, then its value is the address  of the next referenced element, si we set I = value.
		case token.ASTERISK:
			//we set V0 and V1 = value saved from I in memory
			ifx65 := IFX65(2)

			err :=emitter.saveOpcode(ifx65)
			if err != nil{
				return 0,err
			}
			//we set I=value founded previously in I
			iaxy0 := IAXY0(0,1)

			err =emitter.saveOpcode(iaxy0)
			if err != nil{
				return 0,err
			}
			datatype = datatype.(symboltable.Pointer).PointsTo
			emitter.ctxNode = emitter.ctxNode.Children[0]

		}

	}
	return symboltable.GetSize(datatype), nil
}

//saveStackReferenceAddressInI save the address of a reference saved in the stack in I using the register x
//Returns the size of the reference it points to and an error
func (emitter *Emitter)saveStackReferenceAddressInI( x byte, functionCtx *FunctionCtx)(int, error){
	ident := emitter.ctxNode.Value.Literal
	reference, ok := functionCtx.Addresses.References[ident]
	if !ok{
		return 0, errors.New(errorhandler.UnexpectedCompilerError())
	}
	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType)

	//we set I = address position 0 of stack
	iaxy0 := IAXY0(RegisterStackAddres1, RegisterStackAddres2)


	err := emitter.saveOpcode(iaxy0)
	if err != nil{
		return 0, err
	}

	return size, emitter.executeFX1ESafe(x, reference.positionStack)

}

//executeFX1ESafe set an int to vx and then set I = I + vx, if the int is greater than 255 we add vx in a loop
func (emitter *Emitter) executeFX1ESafe(x byte, vx int) error{

	for  vx>255{
		i6xkk := I6XKK(x, 255)
		ifx1e := IFX1E(x)


		err := emitter.saveOpcode(i6xkk)
		if err != nil{
			return err
		}

		err =emitter.saveOpcode(ifx1e)
		if err != nil{
			return err
		}
		vx = vx - 255
	}
	if vx > 0{
		i6xkk := I6XKK(x, byte(vx))
		ifx1e := IFX1E(x)

		err :=emitter.saveOpcode(i6xkk)
		if err != nil{
			return err
		}

		err =emitter.saveOpcode(ifx1e)
		if err != nil{
			return err
		}
	}
	return nil

}

// GetLeafByRight gets the leaf by walking a tree using the right child of each node.
func GetLeafByRight(head *ast.Node) *ast.Node{
	current := head
	for len(current.Children) != 0{
		current = current.Children[len(current.Children)-1]
	}
	return current
}

//saveOpcode save an opcode in the machine code array
func (emitter *Emitter) saveOpcode(opcode Opcode) error{
	emitter.machineCode[emitter.currentAddress] = opcode[0]
	err := emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	emitter.machineCode[emitter.currentAddress] = opcode[1]
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	return nil
}

//obtainSizeParams returns the size of each param of a function
func obtainSizeParams(params []interface{})[]int{
	paramSizes := make([]int,0)
	for _, param := range params{
		paramSizes = append(paramSizes, symboltable.GetSize(param))
	}
	return paramSizes
}