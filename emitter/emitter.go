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
	currentAddress 		   uint16
	offsetStack    		   int //we use this field to know the last address in which we save a variable in the stack
	globalVariables		   map[string]uint16 //we save in globalVariables the address in which each global variable is stored
	scope           	   *symboltable.Scope
	ctxNode         	   *ast.Node
	machineCode    		   [Memory]byte
	translateStatement     map[token.Type]func(*FunctionCtx)error
	translateOperation	   map[token.Type]func(function *FunctionCtx) (*ResultRegIndex, error)
	functions       	   map[string]uint16//we save in functions the address in which each function is stored
	lastIndexSubScope	   int //in the context of a scope, lastIndexSubScope tells the numbers of sub-scopes already written in machineCode

}

func NewEmitter(tree *ast.SyntaxTree, scope *symboltable.Scope)*Emitter{
	emitter := new(Emitter)

	emitter.globalVariables = make(map[string]uint16)
	emitter.functions = make(map[string]uint16)
	emitter.scope = scope
	emitter.lastIndexSubScope = 0
	emitter.ctxNode = tree.Head


	emitter.translateStatement = make(map[token.Type]func(*FunctionCtx)error)

	emitter.translateStatement[token.IF] = emitter._if
	emitter.translateStatement[token.ELSE] = emitter._else
	emitter.translateStatement[token.WHILE] = emitter._while
	emitter.translateStatement[token.EQ] = emitter.assign
	emitter.translateStatement[token.RPAREN] = emitter.voidCall
	emitter.translateStatement[token.RETURN] = emitter._return

	emitter.translateOperation = make(map[token.Type]func(*FunctionCtx)(*ResultRegIndex,error))

	emitter.translateOperation[token.DOLLAR] = emitter.address
	emitter.translateOperation[token.RPAREN] = emitter.parenthesis
	emitter.translateOperation[token.PLUS] = emitter.sum
	emitter.translateOperation[token.MINUS] = emitter.subtraction
	emitter.translateOperation[token.ASTERISK] = emitter.asterisk
	emitter.translateOperation[token.PERCENT] = emitter.mod
	emitter.translateOperation[token.SLASH] = emitter.division
	emitter.translateOperation[token.LTLT] = emitter.shift
	emitter.translateOperation[token.GTGT] = emitter.shift
	emitter.translateOperation[token.OR] = emitter.or
	emitter.translateOperation[token.AND] = emitter.and
	emitter.translateOperation[token.XOR] = emitter.xor
	emitter.translateOperation[token.LAND] = emitter.land
	emitter.translateOperation[token.LOR] = emitter.lor
	emitter.translateOperation[token.EQEQ] = emitter.eqeq
	emitter.translateOperation[token.NOTEQ] = emitter.noteq
	emitter.translateOperation[token.BANG] = emitter.not
	emitter.translateOperation[token.LT] = emitter.ltgt
	emitter.translateOperation[token.GT] = emitter.ltgt
	emitter.translateOperation[token.GTEQ] = emitter.ltgteq
	emitter.translateOperation[token.LTEQ] = emitter.ltgteq
	emitter.translateOperation[token.BOOL] = emitter.boolean
	emitter.translateOperation[token.BYTE] = emitter._byte
	emitter.translateOperation[token.IDENT] = emitter.ident


	emitter.currentAddress = AddressGlobalSection
	return emitter
}


//Start translates the syntax tree into machine code and returns it, return an error if needed
func (emitter *Emitter) Start() ([]byte, error){
	emitter.ctxNode = emitter.ctxNode.Children[0].Children[0] //The tree start with a "" and a EOF node, so we move

	//we save all the global globalVariables into memory
	block := emitter.ctxNode
	for _, child := range block.Children{
		if child.Value.Type == token.LET{
			emitter.ctxNode = child
			err := emitter.globalVariableDeclaration()
			if err != nil{
				return nil, err
			}
		}
	}
	emitter.ctxNode = block
	//we save into memory the primitive functions
	err := emitter.primitiveFunctionsDeclaration()
	if err != nil{
		return nil, err
	}

	mainScope := emitter.scope
	i := 0
	//we save into memory all the functions (including main)
	for _, child := range block.Children{
		if child.Value.Type == token.FUNCTION{
			emitter.scope = mainScope.SubScopes[i]
			i++
			emitter.ctxNode = child
			err = emitter.functionDeclaration()
			if err != nil{
				return nil, err
			}

		}
	}
	emitter.scope = mainScope
	emitter.ctxNode = block


	//The stack section will start in the last available address, which is saved in the vD and vE registers
	vD := byte((emitter.currentAddress & 0xFF00) >> 8)
	VE := byte(emitter.currentAddress & 0x00FF)
	x := byte(RegisterStackAddress1)
	y := byte(RegisterStackAddress2)
	saveV4 := I6XKK(x, vD)
	emitter.machineCode[RomStart] = saveV4[0]
	emitter.machineCode[RomStart+1] = saveV4[1]

	saveV5 := I6XKK(y, VE)
	emitter.machineCode[RomStart+2] = saveV5[0]
	emitter.machineCode[RomStart+3] = saveV5[1]

	//The program will start in the main function, so we jump there
	mainAddress, ok := emitter.functions[token.MAIN]
	if !ok{
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}
	callMain := I2NNN(mainAddress)
	emitter.machineCode[RomStart+4] = callMain[0]
	emitter.machineCode[RomStart+5] = callMain[1]

	//after main we want to repeat
	repeat := I1NNN(RomStart)
	emitter.machineCode[RomStart+6] = repeat[0]
	emitter.machineCode[RomStart+7] = repeat[1]


	return emitter.machineCode[RomStart:Memory-1], nil
}

//function declaration save the instructions of all primitive function in memory
func (emitter *Emitter) primitiveFunctionsDeclaration()error{
	var err error
	err = emitter.drawFontDeclaration()
	if err != nil{
		return err
	}
	/*
	err = emitter.cleanDeclaration()
	if err != nil{
		return err
	}
	err = emitter.setSTDeclaration()
	if err != nil{
		return err
	}
	err = emitter.setDTDeclaration()
	if err != nil{
		return err
	}
	err = emitter.getDTDeclaration()
	if err != nil{
		return err
	}
	err = emitter.randomDeclaration()
	if err != nil{
		return err
	}
	err = emitter.waitKeyDeclaration()
	if err != nil{
		return err
	}
	err = emitter.isKeyPressedDeclaration()
	if err != nil{
		return err
	}
	err = emitter.drawDeclaration()
*/
	return err

}

//drawFontDeclaration save the function drawFont in memory
func (emitter *Emitter) drawFontDeclaration()error{
	emitter.functions[symboltable.FunctionDrawFont] = emitter.currentAddress
	//drawFont has three parameters (a byte in v2, a byte in v3, and byte in v4)
	//it returns a boolean (the value of vf) in v0

	err := emitter.saveOpcode(IFX29(4)) // I = location of sprite for digit V4.
	if err != nil{
		return err
	}

	fontSize :=  byte(5) //every font is represented by 5 bytes

	err = emitter.saveOpcode(IDXYN(2,3,fontSize))
	if err != nil{
		return err
	}

	err = emitter.saveOpcode(I8XY0(0,0xf)) //V0 = Vf
	if err != nil{
		return err
	}

	return emitter.saveOpcode(I00EE())

}

//cleanDeclaration save the function clean in memory
func (emitter *Emitter) cleanDeclaration ()error{
	//clean has not parameters and is a void function that clean the screen
	emitter.functions[symboltable.FunctionClean] = emitter.currentAddress
	err := emitter.saveOpcode(I00E0())
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//setSTDeclaration save the function setST in memory
func (emitter *Emitter) setSTDeclaration ()error{
	emitter.functions[symboltable.FunctionSetST] = emitter.currentAddress
	//setST only has a parameter (a byte) saved in v2, and it is a void function that set sound timer = v2
	err:= emitter.saveOpcode(IFX18(2))
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//setDTDeclaration save the function setDT in memory
func (emitter *Emitter) setDTDeclaration ()error{
	emitter.functions[symboltable.FunctionSetDT] = emitter.currentAddress
	//setST only has a parameter (a byte) saved in v2, and it is a void function that set delay timer = v2
	err :=  emitter.saveOpcode(IFX15(2))
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//getDTDeclaration save the function getDT in memory
func (emitter *Emitter) getDTDeclaration ()error{
	emitter.functions[symboltable.FunctionGetDT] = emitter.currentAddress
	//getDT has no parameters and it return a byte (the value of delay timer) in v0
	err := emitter.saveOpcode(IFX07(0))
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//randomDeclaration save the function random in memory
func (emitter *Emitter) randomDeclaration ()error{
	emitter.functions[symboltable.FunctionRandom] = emitter.currentAddress
	//random has no parameters and it returns a random byte (in v0)
	err := emitter.saveOpcode(ICXKK(0, 0xFF))
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//waitKeyDeclaration save the function waitKey in memory
func (emitter *Emitter) waitKeyDeclaration ()error{
	emitter.functions[symboltable.FunctionWaintKet] = emitter.currentAddress
	//waitKey has no parameters and it returns the value of a key pressed in v0
	err := emitter.saveOpcode(IFX0A(0))
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}

//isKeyPressedDeclaration save the function isKeyPressed in memory
func (emitter *Emitter) isKeyPressedDeclaration ()error{
	emitter.functions[symboltable.FunctionIsKeyPressed] = emitter.currentAddress
	//isKeyPressed has one parameter in v2(a byte) and it returns a bool in v0
	err := emitter.saveOpcode(I6XKK(1, True)) //V1 = True

	if err != nil{
		return err
	}
	err = emitter.saveOpcode(IEX9E(2)) //If the key saved in v2 was pressed we skip the next instruction
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(I6XKK(1, False)) //If the key saved in v2 was not pressed we set v1 = False
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(I8XY0(0,1)) //V0=V1
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}
//drawDeclaration save the function draw in memory
func (emitter *Emitter) drawDeclaration ()error {
	emitter.functions[symboltable.FunctionDraw] = emitter.currentAddress
	//draw has four parameters (a byte in v2, a byte in v3, a byte in v4, and pointer in v5 and v6)
	//it returns a boolean (the value of vf) in v0


	dxynAddress := emitter.functions[symboltable.FunctionDraw]+6 //address in which we want dynamically write the opcode

	err := emitter.saveOpcode(I6XKK(0, 0xD2)) //v0=0xD2 (v0 =0xDX)
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(I6XKK(1, 0X30)) //v1=0x30 (v1 = 0xY0)
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(I8XY1(1, 4)) //v1=v1 | v4, (v1 = 0xYN)
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(IANNN(dxynAddress)) //I =dxynAddress
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(IFX55(1)) //save v0 and v1 in dxynAddress (writing the opcode)
	if err != nil{
		return err
	}
	err = emitter.saveOpcode(I9XY1(5,6)) //I=Pointer
	if err != nil{
		return err
	}

	//the next 2 bytes in memory are for the DXYN opcode that was just dynamically generated

	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	err = emitter.moveCurrentAddress()
	if err != nil{
		return err
	}
	return emitter.saveOpcode(I00EE())

}
//function declaration save all the instructions of a function in memory
func (emitter *Emitter) functionDeclaration()error{
	const IDENT = 0
	const ARG = 1
	const BLOCK = 3
	emitter.lastIndexSubScope = 0

	//we backup the offsetStack so we can update it after compiling the function
	offsetBackup := emitter.offsetStack

	//we store the address in which the function is saved in the map of functions
	functionName := emitter.ctxNode.Children[IDENT].Value.Literal
	emitter.functions[functionName] = emitter.currentAddress //the function starts at the current address

	mainScope := emitter.scope
//	emitter.scope = emitter.scope.SubScopes[len(emitter.functions)-1]
	ctxReferences := NewStackReferences()

	//we save the arguments in the stack
	fn  := emitter.ctxNode
	params := make([]string, 0)
	hasParams := false
	if len(emitter.ctxNode.Children[ARG].Children)>0{
		hasParams = true
		funcSymbol := emitter.scope.Symbols[functionName]
 		sizeParams := obtainSizeParams(funcSymbol.DataType.(symboltable.Function).Args)
		i := 0
		var err error
		emitter.ctxNode = emitter.ctxNode.Children[ARG].Children[0]
		for emitter.ctxNode.Value.Type == token.COMMA{
			comma := emitter.ctxNode
			emitter.ctxNode = emitter.ctxNode.Children[0]
			err = emitter.saveParamsInStack(&params, ctxReferences, i, sizeParams)
			if err != nil {
				return err
			}
			i++
			emitter.ctxNode = comma
			emitter.ctxNode = emitter.ctxNode.Children[1]
		}
		err = emitter.saveParamsInStack(&params, ctxReferences, i, sizeParams)

		if err != nil {
			return err
		}
	}

	emitter.ctxNode = fn

	//we declare all variables in the stack
	emitter.ctxNode = emitter.ctxNode.Children[BLOCK]
	err := emitter.declareInStack(ctxReferences)
	if err != nil {
		return err
	}

	emitter.ctxNode = fn


	ctxFunction := NewCtxFunction(ctxReferences)
	if hasParams{
		emitter.scope = emitter.scope.SubScopes[0]
	}
	//we write the rest of the statements in memory
	for _, child := range fn.Children[BLOCK].Children {

		emitter.ctxNode = child
		//we jump let stmts because we already save them
		translateStmt, ok := emitter.translateStatement[emitter.ctxNode.Value.Type]
		if ok {
			err := translateStmt(ctxFunction)
			if err != nil {
				return err
			}
		}

	}

	emitter.scope = mainScope
	emitter.offsetStack = offsetBackup
	return nil
}

//saveParamsInStack declare params in the stack and save its values there, returns the an error if needed
func (emitter *Emitter) saveParamsInStack(params *[]string, ctxReferences *Stack, i int, sizeParams []int) error {
	const IDENT = 0
	//first we declare them in the stack
	paramIdent := emitter.ctxNode.Children[IDENT].Value.Literal
	*params = append(*params, paramIdent)
	err := emitter.let(ctxReferences)
	if err != nil {
		return  err
	}
	//then we set its value
	//the argument i is the register i + 2 (we use v0 and v1 to operate)

	err = emitter.saveOpcode(I8XY0(0, byte(i+2))) //v0 = v(i+2)
	if err != nil {
		return  err
	}
	if sizeParams[i] == 2 {
		  //v1 = v(i+2)
		err = emitter.saveOpcode(I8XY0(1, byte(i+3)))
		if err != nil {
			return err
		}

	}

	err = emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = stack address
	if err != nil {
		return err
	}
	reference, _ := ctxReferences.GetReference(paramIdent)

	err = emitter.saveFX1ESafely(2, reference.positionInStack) // I = I + V2
	if err != nil {
		return err
	}

	err = emitter.saveOpcode(IFX55(byte(sizeParams[i]-1))) //we save v0 (or v0 and v1) in memory
	if err != nil {
		return err
	}

	emitter.offsetStack = emitter.offsetStack + sizeParams[i]
	return nil
}

//globalVariableDeclaration assigns an address to a global variable and updates the current address.
func (emitter *Emitter) globalVariableDeclaration() error{
	let := emitter.ctxNode
	ident := let.Children[0].Value.Literal
	symbol, ok := emitter.scope.Symbols[ident]
	if !ok{
		return errors.New(errorhandler.UnexpectedCompilerError())
	}

	size := symboltable.GetSize(symbol.DataType)
	emitter.globalVariables[ident] = emitter.currentAddress

	for i := 0; i < size; i++{
		emitter.machineCode[emitter.currentAddress] = 0
		err := emitter.moveCurrentAddress()
		if err != nil{
			return err
		}
	}
	return nil


}

//moveCurrentAddress moves the current address by one, and if it's out of bounds of the memory it return a error
func (emitter *Emitter) moveCurrentAddress() error{
	emitter.currentAddress++
	if emitter.currentAddress > Memory{
		return errors.New(errorhandler.NotEnoughMemory())
	}else {
		return nil
	}

}

//declareInStack saves all variables of a function in its stack
func (emitter *Emitter) declareInStack(ctxReferences *Stack) error {
	backupCtxNode := emitter.ctxNode
	backupScope := emitter.scope

	iSubScope := 0

	for _, child := range emitter.ctxNode.Children{
		switch child.Value.Type {
		case token.WHILE:
			emitter.ctxNode = child.Children[1]
			emitter.scope = emitter.scope.SubScopes[iSubScope]
			ctxReferences.AddSubReferences()

			err := emitter.declareInStack(ctxReferences.SubReferences[iSubScope])
			iSubScope++

			emitter.ctxNode = backupCtxNode
			emitter.scope = backupScope
			if err != nil{
					return err
			}

		case token.IF:
			emitter.ctxNode = child.Children[1]
			emitter.scope = emitter.scope.SubScopes[iSubScope]
			ctxReferences.AddSubReferences()
			err := emitter.declareInStack(ctxReferences.SubReferences[iSubScope])
			iSubScope++
			emitter.ctxNode = backupCtxNode
			emitter.scope = backupScope
			if err != nil{
				return err
			}
		case token.ELSE:
			for j:=0; j<2;j++{
				emitter.ctxNode = child.Children[j+1]
				emitter.scope = emitter.scope.SubScopes[iSubScope]
				ctxReferences.AddSubReferences()
				err := emitter.declareInStack(ctxReferences.SubReferences[iSubScope])
				iSubScope++
				emitter.ctxNode = backupCtxNode
				emitter.scope = backupScope
				if err != nil{
					return err
				}
			}
		case token.LET:
			backup := emitter.ctxNode
			emitter.ctxNode = child
			err := emitter.let(ctxReferences)
			if err != nil{
				return nil
			}
			emitter.ctxNode = backup
		}
	}

	emitter.ctxNode = backupCtxNode
	return nil

}

//let saves a specific variable in the stack of a function and update ctxReferences
func (emitter *Emitter) let(ctxReferences *Stack) error{
	IDENT := 0
	ident := emitter.ctxNode.Children[IDENT].Value.Literal
	ctxReferences.AddReference(ident, emitter.offsetStack)
	symbol, ok := emitter.scope.Symbols[ident]
	if !ok{
		return errors.New(errorhandler.UnexpectedCompilerError())
	}
	size := symboltable.GetSize(symbol.DataType)

	for size > 16{
		err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) //I = stack
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(I6XKK(0, byte(emitter.offsetStack)))                 //v0 = offset
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(IFX1E(0))                  // I = I + V0
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(IFX55(15))
		size -=16
		emitter.offsetStack += 16
	}
	if size > 0{
		err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) //I = stack
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(I6XKK(0, byte(emitter.offsetStack)))                 //v0 = offset
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(IFX1E(0))                  // I = I + V0
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(IFX55(byte(size-1)))
		emitter.offsetStack += size
	}

	return nil
}


//assign translates the assign statement to opcodes and write it in emitter.machineCode
func (emitter *Emitter)assign(functionCtx *FunctionCtx)error {
	const SAVEIN = 0
	const TOSAVE = 1

	assignBackup := emitter.ctxNode
	emitter.ctxNode = emitter.ctxNode.Children[TOSAVE]
	//depending on the data type of the right side of the assign translateOperation save it in v0 or v0 and v1
	size, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil {
		return err
	}
	emitter.ctxNode = assignBackup
	emitter.ctxNode = emitter.ctxNode.Children[SAVEIN]

	//if in the right side of the assign there is only a identifier, then it is a reference to a variable
	if emitter.ctxNode.Value.Type == token.IDENT {
		ident := emitter.ctxNode.Value.Literal
		_, referenceExists := functionCtx.stack.GetReference(ident) //we check if it is saved in the stack
		if !referenceExists{                                        //if it is not saved in the stack it must be a global variable

			_, isAGlobalVariable := emitter.globalVariables[ident]
			if !isAGlobalVariable{
				return errors.New(errorhandler.UnexpectedCompilerError())
			}else {
				//if the reference is to a global variable we update I = address
				_, err := emitter.saveGlobalReferenceAddressInI(2, 3)
				if err != nil{
					return err
				}
				//and then we store registers V0 through Vsize in memory starting at location I

				err =emitter.saveOpcode(IFX55(byte(size-1)))

				if err != nil{
					return err
				}
				return nil

			}
		}else{//if it is stored in the stack
			// we look for the address of the reference in the stack, and we set I = address

			_,err := emitter.saveStackReferenceAddressInI(2, functionCtx)
			if err != nil{
				return err
			}
			//then we save v0 (and maybe v1) there
			err =emitter.saveOpcode(IFX55(byte(size-1)))

			if err != nil{
				return err
			}
			return nil


		}


	}else{ //if in the right side of the assign there is a sequence of characters before a ident, it's a dereference

		//we backup v0 (and maybe v1) in v2 (and maybe v3) because we need v0 and v1 to operate
		err = emitter.saveOpcode(I8XY0(2, 0))
		if err!=nil{
			return err
		}
		if symboltable.GetSize(size) > 1{

			err =emitter.saveOpcode(I8XY0(3, 1))

			if err!=nil{
				return err
			}
		}
		//we search the address of the dereference and we save it in I
		_, err = emitter.saveDereferenceAddressInI(functionCtx)
		if err != nil{
			return err
		}
		//to save a vx (and vy) in memory we need them in v0 (and v1), so we save them there again

		err =emitter.saveOpcode( I8XY0(0, 2))

		if err!=nil{
			return err
		}
		if symboltable.GetSize(size) > 1{
			err =emitter.saveOpcode(I8XY0(1, 3))

			if err!=nil{
				return err
			}
		}
		err =emitter.saveOpcode(IFX55(byte(size-1)))

		if err != nil{
			return err
		}

		return nil

	}

}

//_return save in v0 the value a function returns
func (emitter *Emitter) _return(functionCtx *FunctionCtx) error {
	if len(emitter.ctxNode.Children) != 0{
		returnBackup := emitter.ctxNode
		emitter.ctxNode = emitter.ctxNode.Children[0]
		_, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
		if err != nil{
			return  err
		}
		emitter.ctxNode = returnBackup
	}
	return emitter.saveOpcode(I00EE())

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
	err = emitter.saveOpcode(I3XKK(0, True))  //if v0 = true we skip the next instruction
	if err != nil{
		return err
	}
	//the next instruction is a jump to the memory address after the block
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
	err = emitter.saveOpcode(I3XKK(0, True))  //if v0 = true we skip the next instruction
	if err != nil{
		return err
	}
	//the next instruction is a jump to the memory address after the if block
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
	i1nnn = I1NNN(emitter.currentAddress)

	emitter.machineCode[lineAfterIf] = i1nnn[0]
	emitter.machineCode[lineAfterIf+1] = i1nnn[1]
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
	err = emitter.saveOpcode(I3XKK(0, True)) //if v0 = true we skip the next instruction
	if err != nil{
		return err
	}
	//the next instruction is a jump to the memory address after the while
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
	addressesBackup := functionCtx.stack
	emitter.scope = emitter.scope.SubScopes[emitter.lastIndexSubScope]
	functionCtx.stack = functionCtx.stack.SubReferences[emitter.lastIndexSubScope]
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
	functionCtx.stack = addressesBackup
	return nil
}

//parenthesis analyze the context of a parenthesis and delegate the operation,
//returns the index of registers in which the result of the operation was stored and an error if needed
func (emitter *Emitter)parenthesis(functionCtx *FunctionCtx)(*ResultRegIndex, error){
	if emitter.ctxNode.Children[0].Value.Type == token.IDENT{
		return emitter.call(functionCtx)
	}
	emitter.ctxNode = emitter.ctxNode.Children[0] //skip node
	return emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
}

//voidCall translates a void call to opcodes and write it in emitter.machineCode
func (emitter *Emitter)voidCall(functionCtx *FunctionCtx)error{
	_,err := emitter.call(functionCtx)
	return err
}

//call translates a function call to opcodes and write it in emitter.machineCode,
//returns the indexes of registers in which the return value is stored and an error if needed
func (emitter *Emitter)call(functionCtx *FunctionCtx)(*ResultRegIndex, error) {
	const IDENT = 0

	//we first backup all registers of the current function in  the stack
	err := emitter.backupRegistersInMemory(functionCtx)
	if err != nil {
		return nil, err
	}
	emitter.offsetStack+=AmountOfRegistersToOperate
	//then we save the params of the function call in registers
	err = emitter.saveParamsInRegisters(functionCtx)
	if err != nil {
		return nil, err
	}

	ident := emitter.ctxNode.Children[IDENT].Value.Literal
	//we call the function
	fnAddress, _ := emitter.functions[ident]
	err = emitter.saveOpcode(I2NNN(fnAddress))
	if err != nil {
		return nil, err
	}
	//if it was not a void function, now the return value is in vf.

	//if the function we call was not a void function, then we save in memory a backup of the return value,
	//because we will need the registers v0
	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType.(symboltable.Function).Return)
	if size != 0 {
		err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = stack address
		if err != nil {
			return nil, err
		}
		err = emitter.saveFX1ESafely(2, emitter.offsetStack) //I = I + offset
		if err != nil {
			return nil, err
		}
		emitter.offsetStack += size
		err = emitter.saveOpcode(IFX55(byte(size-1))) //TODO: By now size is always 1
		if err != nil {
			return nil, err
		}
	}

	//then we save again the previous registers in memory
	err = emitter.takeRegistersFromMemory(functionCtx)
	if err != nil {
		return nil, err
	}

	//and if it wasn't a void function, we save again the return values in a register

	if size != 0 {
		regIndex, ok := functionCtx.registerHandler.AllocSimple()
		if !ok{
			line := emitter.ctxNode.Value.Line
			err := errors.New(errorhandler.TooManyRegisters(line))
			return nil, err
		}
		err = emitter.saveFX1ESafely(regIndex.lowBitsIndex, emitter.offsetStack-size) //I = End of the register backup section/start of the return value backup section
		if err != nil {
			return nil, err
		}
		err = emitter.saveOpcode(IFX65(byte(size-1)))
		if err != nil {
			return nil, err
		}
		//now the return value is again in 0
		err = emitter.saveOpcode(I8XY0(regIndex.lowBitsIndex, 0)) //Vx = V0
		if err != nil {
			return nil, err
		}


	}
	return nil, nil

}

func (emitter *Emitter) takeRegistersFromMemory(functionCtx *FunctionCtx)  error {
	err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = address stack
	if err != nil {
		return err
	}
	err = emitter.saveFX1ESafely(0, emitter.offsetStack) //I = I + offset
	if err != nil {
		return  err
	}
	err = emitter.saveOpcode(IFX65(AmountOfRegistersToOperate-1))
	if err != nil {
		return  err
	}
	return nil
}

func (emitter *Emitter) backupRegistersInMemory(functionCtx *FunctionCtx) error{
	err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = address stack
	if err != nil {
		return err
	}
	err = emitter.saveFX1ESafely(0, emitter.offsetStack) //I = I + offset
	if err != nil {
		return  err
	}
	err = emitter.saveOpcode(IFX55(AmountOfRegistersToOperate-1))
	if err != nil {
		return  err
	}

	return nil
}

func (emitter *Emitter) saveParamsInRegisters(functionCtx *FunctionCtx) error {
	const PARAMS = 1
	backupNode := emitter.ctxNode

	backupRegisterHandler := functionCtx.registerHandler
	functionCtx.registerHandler := NewRegisterHandler()
	/

	paramSizes := 0
	i:=2
	if len(emitter.ctxNode.Children) > 1 { //we ask if it has any param

		emitter.ctxNode = emitter.ctxNode.Children[PARAMS]
		for emitter.ctxNode.Value.Type == token.COMMA {
			backupComma := emitter.ctxNode
			emitter.ctxNode = emitter.ctxNode.Children[0]
			functionCtx.registerHandler.reserveRegister(byte(i))
			//we save in v0(and maybe v1) the value of the parameter being analyzed
			resultRegIndex, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
			if err != nil {
				return  err
			}
			err = emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = address stack
			if err != nil {
				return  err
			}
			err = emitter.saveFX1ESafely(2, emitter.offsetStack) //we move I = last address in the stack
			if err != nil {
				return err
			}

			err = emitter.saveOpcode(IFX55(byte(size-1))) //we save the param in the stack
			if err != nil {
				return err
			}

			paramSizes += size
			emitter.offsetStack += size
			emitter.ctxNode = backupComma
			emitter.ctxNode = emitter.ctxNode.Children[1]
		}
		//emitter.ctxNode = emitter.ctxNode.Children[0]
		//we save in v0(and maybe v1) the value of the parameter being analyzed
		size, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
		if err != nil {
			return  err
		}
		err = emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = address stack
		if err != nil {
			return err
		}
		err = emitter.saveFX1ESafely(2, emitter.offsetStack) //we move I = last address in the stack
		if err != nil {
			return  err
		}

		paramSizes += size
		err = emitter.saveOpcode(IFX55(byte(size-1))) //we save the param in the stack
		if err != nil {
			return err
		}

	}
	//now that all the params are saved in the stack, we store them in registers

	emitter.offsetStack = backupStack

	if emitter.offsetStack - 2 < 0{
		toSubtract := (emitter.offsetStack-2)*(-1)

		// we save VD and VE in V0 and V1 and the number to subtract in v2
		err := emitter.saveOpcode(I8XY0(0,0xD))
		if err != nil{
			return err
		}

		err = emitter.saveOpcode(I8XY0(1,0xE))
		if err != nil{
			return err
		}

		err = emitter.saveOpcode(I6XKK(2,byte(toSubtract)))
		if err != nil{
			return err
		}

		//we first subtract v1 = v1 - v2
		err = emitter.saveOpcode(I8XY5(1,2))
		if err != nil{
			return err
		}
		//because we already use v2, we can now use it as an aux, v2 = 1
		err = emitter.saveOpcode(I6XKK(2,1))
		if err != nil{
			return err
		}
		//if vf = false, then v1 - v2 < 0, so we need to set v0 = v0 - 1
		err = emitter.saveOpcode(I4XKK(0xf,False))
		if err != nil{
			return err
		}
		err = emitter.saveOpcode(I8XY5(0,2))
		if err != nil{
			return err
		}

		//now we set I = stack - offset - 2
		err = emitter.saveOpcode(I9XY1(0, 1))
		if err != nil {
			return  err
		}

	}else{
		err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2)) // I = address stack
		if err != nil {
			return  err
		}
		err = emitter.saveFX1ESafely(2, emitter.offsetStack-2) //we move I = I + (offsetStack-2) because we want the params
		//to start in v2
		if err != nil {
			return  err
		}

	}

	err := emitter.saveOpcode(IFX65(byte(paramSizes -1 + 2))) //we read the param in the stack (from v2 through v(ParamSize+2))
	if err != nil {
		return  err
	}
	emitter.ctxNode = backupNode
	return nil
}

//_byte save a byte in a registers. Return the register index in which the byte was stored and an error if needed
func (emitter *Emitter) _byte(functionCtx *FunctionCtx) (*ResultRegIndex, error) {
	regIndex, ok := functionCtx.registerHandler.AllocSimple()
	if !ok{
		line := emitter.ctxNode.Value.Line
		err := errors.New(errorhandler.TooManyRegisters(line))
		return regIndex, err
	}

	kk, _ := strconv.Atoi(emitter.ctxNode.Value.Literal)
	err := emitter.saveOpcode(I6XKK(regIndex.lowBitsIndex, byte(kk))) // Vx = Byte
	if err != nil{
		return regIndex, err
	}

	return regIndex, nil
}

//boolean save a bool in a registers. Return the register index in which the bool was stored and an error if needed
func (emitter *Emitter)boolean(functionCtx *FunctionCtx) (*ResultRegIndex, error) {
	regIndex, ok := emitter.functionCtx.AllocSimple()
	if !ok{
		line := emitter.ctxNode.Value.Line
		err := errors.New(errorhandler.TooManyRegisters(line))
		return nil, err
	}

	var kk byte
	if emitter.ctxNode.Value.Literal == token.TRUE{
		kk = True
	}else{
		kk = False
	}
	err := emitter.saveOpcode(I6XKK(regIndex.lowBitsIndex, kk))
	if err != nil{
		return nil, err
	}
	return regIndex, nil
}

//ltgt translates < and > to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter)ltgt(functionCtx *FunctionCtx) (int, error){
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}

	// if our operands are simples we just need to check if v0 is lesser/greater than v2.
	//we first save this information in vf and then we set v0 = vf
	if sizeOperands[0] == 1{
		switch emitter.ctxNode.Value.Type {
		case token.GT:
			err = emitter.saveOpcode(I8XY5(0,2))
			if err != nil{
				return 0, err
			}
		case token.LT:
			err = emitter.saveOpcode(I8XY7(0,2))
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())
		}
		err = emitter.saveOpcode(I8XY0(0,0xf))
		if err != nil{
			return 0, err
		}
		return 1, nil
	}else{
		//if we are comparing pointers we first compare v0 with v2
		switch emitter.ctxNode.Value.Type {
		case token.GT:
			err = emitter.saveOpcode(I8XY5(0,2))
			if err != nil{
				return 0, err
			}
		case token.LT:
			err = emitter.saveOpcode(I8XY7(0,2))
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())

		}
		err = emitter.saveOpcode(I3XKK(0xf,0)) //if vf = 0 we keep analyzing
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I1NNN(emitter.currentAddress + 5)) //if vf = 1 we know that the result is true and jump to the end
		if err != nil{
			return 0, err
		}
		//if vf = 0 we ask if v0 == v2 with a xor
		err = emitter.saveOpcode(I8XY3(0,2))
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I4XKK(0,0)) //if v0 != v2 we skip the next opcode
		if err != nil{
			return 0, err
		}


		//if v0 == v2 we need to analyze v1 and v3
		switch emitter.ctxNode.Value.Type {
		case token.GT:
			err = emitter.saveOpcode(I8XY5(0,2))
			if err != nil{
				return 0, err
			}
		case token.LT:
			err = emitter.saveOpcode(I8XY7(0,2))
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())

		}
		//then we save the result in v0
		err = emitter.saveOpcode(I8XY0(0,0xf))
		if err != nil{
			return 0, err
		}
		return 2, nil
	}

}

//ltgteq translates <= and >= to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter)ltgteq(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil {
		return 0, err
	}

	if sizeOperands[0] == 1{
		err = emitter.saveOpcode(I6XKK(0xf,True)) //vf = 1
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I8XY3(0,2)) //we ask v0 == v2 with a xor
		if err != nil{
			return 0, err
		}

		err = emitter.saveOpcode(I3XKK(0,0)) //if v0 = 0 then v0 was equal to v2 and we skip the next opcode
		if err != nil{
			return 0, err
		}
		switch emitter.ctxNode.Value.Type {
		case token.GTEQ:
			err = emitter.saveOpcode(I8XY5(0,2))//if v0 wasn't equal to v2, we ask if v0 > v2 and store the result in vf
			if err != nil{
				return 0, err
			}
		case token.LTEQ:
			err = emitter.saveOpcode(I8XY7(0,2)) //if v0 wasn't equal to v2, we ask if v0 < v2 and store the result in vf
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())

		}
		//then we save the result in v0
		err = emitter.saveOpcode(I8XY0(0,0xf))
		if err != nil{
			return 0, err
		}

		return 1, nil
	}else{
		//first we ask if v0 is greater/lesser than v2 and store the result in vf
		switch emitter.ctxNode.Value.Type {
		case token.GTEQ:
			err = emitter.saveOpcode(I8XY5(0,2))
			if err != nil{
				return 0, err
			}
		case token.LTEQ:
			err = emitter.saveOpcode(I8XY7(0,2))
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())

		}

		err = emitter.saveOpcode(I3XKK(0xf,0)) //if vf = 0 we skip  the next opcode
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I1NNN(emitter.currentAddress+9)) //if vf = 1 we skip  7 opcodes because we know the result is true
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I8XY3(0,2)) //we ask v0 == v2 with a xor. v0 = 0 if ture
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I3XKK(0,0)) //if v0 = 0 we skip  the next opcode
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I1NNN(emitter.currentAddress+6)) //if v0 != 0 we skip 4 opcodes because we know the result is false
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I8XY3(0,2)) //we ask v1 == v3 with a xor. v1 = 0 if ture
		if err != nil{
			return 0, err
		}
		// if v1 = 0 we set vf = 1 and jump to the end
		err = emitter.saveOpcode(I4XKK(1,0))
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I6XKK(0xf,True)) //vf = 1
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I3XKK(1,0))
		if err != nil{
			return 0, err
		}
		//first we ask if v1 is greater/lesser than v3 and store the result in vf
		switch emitter.ctxNode.Value.Type {
		case token.GTEQ:
			err = emitter.saveOpcode(I8XY5(0,2))
			if err != nil{
				return 0, err
			}
		case token.LTEQ:
			err = emitter.saveOpcode(I8XY7(0,2))
			if err != nil{
				return 0, err
			}
		default:
			return 0, errors.New(errorhandler.UnexpectedCompilerError())

		}
		//we set v0 = vf
		err = emitter.saveOpcode(I8XY0(0,0xf))
		if err != nil{
			return 0, err
		}


		return 2, nil
	}
}
//noteq translates a != to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) noteq(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil {
		return 0, err
	}
	//if the operands are simple data types we do a xor between v0 and v2,
	//if they are equal v0 = 0
	err = emitter.saveOpcode(I8XY3(0,2))
	if err != nil{
		return 0, err
	}

	if sizeOperands[0] == 2 {
		//if not, we do v0 = v0 ^ v2, v1 = v1 ^ v3, v0 = v0 | v1
		err = emitter.saveOpcode(I8XY3(1,3))
		if err != nil{
			return 0, err
		}
		err = emitter.saveOpcode(I8XY1(0,1))
		if err != nil{
			return 0, err
		}
	}
	return 1, nil

}
//eqeq translates a == to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) eqeq(functionCtx *FunctionCtx) (int, error) {
	//we do the same than in !=, but with a not at the end
	_, err := emitter.noteq(functionCtx)
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I6XKK(1,True))
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I8XY3(0,1))
	if err != nil{
		return 0, err
	}

	return 1, nil

}

//not translates a ! to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) not(functionCtx *FunctionCtx) (int, error) {
	emitter.ctxNode = emitter.ctxNode.Children[0]
	_, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return 0, err
	}
	//we set v0 = v0 ^ true
	err = emitter.saveOpcode(I6XKK(1,True))
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I8XY3(0,1))
	if err != nil{
		return 0, err
	}

	return 1, nil

}


//land translates a && to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) land(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//V0 = V0 & V2
	err = emitter.saveOpcode(I8XY2(0,2))
	if err != nil{
		return 0, err
	}

	return 1, nil
}

//lor translates a || to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) lor(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//V0 = V0 | V2
	err = emitter.saveOpcode(I8XY1(0,2))
	if err != nil{
		return 0, err
	}

	return 1, nil
}

//or translates a | to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) or(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//V0 = V0 | V2
	err = emitter.saveOpcode(I8XY1(0,2))
	if err != nil{
		return 0, err
	}

	if sizeOperands[0] > 1{

		//V1 = V1 | V3
		err = emitter.saveOpcode(I8XY1(1,3))
		if err != nil{
			return 0, err
		}
	}
	return sizeOperands[0], nil
}

//and translates a & to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) and(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//V0 = V0 & V2
	err = emitter.saveOpcode(I8XY2(0,2))
	if err != nil{
		return 0, err
	}

	if sizeOperands[0] > 1{

		//V1 = V1 & V3
		err = emitter.saveOpcode(I8XY2(1,3))
		if err != nil{
			return 0, err
		}
	}
	return sizeOperands[0], nil
}

//xor translates a ^ to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) xor(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//V0 = V0 ^ V2
	err = emitter.saveOpcode(I8XY3(0,2))
	if err != nil{
		return 0, err
	}

	if sizeOperands[0] > 1{

		//V1 = V1 ^  V3
		err = emitter.saveOpcode(I8XY3(1,3))
		if err != nil{
			return 0, err
		}
	}
	return sizeOperands[0], nil
}


//sum translates a sum to opcodes and write it in emitter.machineCode,
//returns the indexes of registers in which is stored the result and an error
func (emitter *Emitter) sum(functionCtx *FunctionCtx) (*ResultRegIndex, error) {
	leftRegIndex, rightRegIndex, err := emitter.solveOperands(functionCtx)
	if err != nil {
		return nil, err
	}

	if leftRegIndex.isPointer{
		//the sum has the same size than the left operand
		sumRegIndex, ok := functionCtx.registerHandler.AllocPointer()
		if !ok{
			line := emitter.ctxNode.Value.Line
			err := errors.New(errorhandler.TooManyRegisters(line))
			return nil, err
		}
		//if the left operands is a pointer we first sum vLeft1 = vLeft1 + vRight
		err = emitter.saveOpcode(I8XY4(leftRegIndex.lowBitsIndex, rightRegIndex.lowBitsIndex))
		if err != nil {
			return nil, err
		}
		//if carry = true, then vLeft1 + vRight > 255, so we need to set vLeft0 = vLeft0 + 1
		err = emitter.saveOpcode(I4XKK(Carry, True))
		if err != nil {
			return nil, err
		}
		err = emitter.saveOpcode(I7XKK(leftRegIndex.highBitsIndex, 1))
		if err != nil {
			return nil, err
		}
		//we save viSum0 = vLeft0, viSum1 = vLeft1
		err = emitter.saveOpcode(I8XY0(sumRegIndex.highBitsIndex, leftRegIndex.highBitsIndex))
		if err != nil {
			return nil, err
		}

		err = emitter.saveOpcode(I8XY0(sumRegIndex.lowBitsIndex, leftRegIndex.lowBitsIndex))
		if err != nil {
			return nil, err
		}
		functionCtx.registerHandler.Free(leftRegIndex)
		functionCtx.registerHandler.Free(rightRegIndex)

		return sumRegIndex, nil
	}else{
		sumRegIndex, ok := functionCtx.registerHandler.AllocSimple()
		if !ok{
			line := emitter.ctxNode.Value.Line
			err := errors.New(errorhandler.TooManyRegisters(line))
			return nil, err
		}
		//if the left operand is a simple we just sum vLeft = vLeft +vRight, and we return the result in
		//a single register v_isum
		err := emitter.saveOpcode(I8XY4(leftRegIndex.lowBitsIndex, rightRegIndex.lowBitsIndex))
		if err != nil {
			return nil, err
		}
		err = emitter.saveOpcode(I8XY0(sumRegIndex.lowBitsIndex, leftRegIndex.lowBitsIndex))
		if err != nil {
			return nil, err
		}

		functionCtx.registerHandler.Free(leftRegIndex)
		functionCtx.registerHandler.Free(rightRegIndex)
		return sumRegIndex, nil

	}


}


//subtraction translates a subtraction to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) subtraction(functionCtx *FunctionCtx) (int, error) {
	sizeOperands, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}
	//if the left operand is a simple data type we just subtract v0 = v0 - v2
	if sizeOperands[0] == 1{
		err := emitter.saveOpcode(I8XY5(0,2))
		if err != nil{
			return 0, err
		}
		return 1, nil
	}
	//if the left operands is a pointer we first subtract v1 = v1 - v2
	err = emitter.saveOpcode(I8XY5(1,2))
	if err != nil{
		return 0, err
	}
	//because we already use v2, we can now use it as an aux, v2 = 1
	err = emitter.saveOpcode(I6XKK(2,1))
	if err != nil{
		return 0, err
	}
	//if vf = false, then v1 - v2 < 0, so we need to v0 = v0 - 1
	err = emitter.saveOpcode(I4XKK(0xf,False))
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I8XY5(0,2))
	if err != nil{
		return 0, err
	}
	return 2, nil
}

//shift translates a subtraction to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) shift(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil {
		return 0, err
	}
	//we set v3= 1 to use it as an aux
	err = emitter.saveOpcode(I6XKK(3,1))
	if err != nil{
		return 0, err
	}

	//we shift v0 by 1
	switch emitter.ctxNode.Value.Type {
	case token.GTGT:
		err = emitter.saveOpcode(I8XY6(0))
		if err != nil{
			return 0, err
		}
	case token.LTLT:
		err = emitter.saveOpcode(I8XYE(0))
		if err != nil{
			return 0, err
		}
	default:
		return 0, errors.New(errorhandler.UnexpectedCompilerError())

	}

	//v2 = v2 - 1
	err = emitter.saveOpcode(I8XY5(2, 3))
	if err != nil{
		return 0, err
	}

	//if v2 != 0 we keep shifting
	err = emitter.saveOpcode(I3XKK(0, 0))
	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I1NNN(emitter.currentAddress-3))
	if err != nil{
		return 0, err
	}

	return 1, nil

}



//multiplication translates a multiplication to opcodes and write it in emitter.machineCode,
//returns the size of the datatype of the result and an error
func (emitter *Emitter) multiplication(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode( I4XKK(0, 0) )//if v0 != 0 we skip the next opcode

	if err != nil{
		return 0, err
	}
	//if v0 =0, the result is 0 and we skip the operation
	skipMultiplication := I1NNN(emitter.currentAddress+6)
	err = emitter.saveOpcode(skipMultiplication)
	if err != nil{
		return 0, err
	}
	//v1 = 1
	err = emitter.saveOpcode(I6XKK(1,1))
	if err != nil{
		return 0, err
	}

	//v2 = v2 + v2
	err = emitter.saveOpcode(I8XY4(2,2))
	if err != nil{
		return 0, err
	}

	//v2 = v2 - v1
	err = emitter.saveOpcode(I8XY5(2,1))
	if err != nil{
		return 0, err
	}

	//if v2 = 0 we skip the next opcode
	err = emitter.saveOpcode(I3XKK(2,0))
	if err != nil{
		return 0, err
	}

	//if v2 != 0 we keep iterating the loop
	err = emitter.saveOpcode(I1NNN(emitter.currentAddress-3))
	if err != nil{
		return 0, err
	}

	return 1, nil
}

//mod translates a % to opcodes and write it in emitter.machineCode, return the size of the datatype of the result and an error
func (emitter *Emitter) mod(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}


	err = emitter.saveOpcode(I4XKK(0, 0))//if v0 != 0 we skip the next opcode

	if err != nil{
		return 0, err
	}
	//if v0 =0, the result is 0 and we skip the operation
	skipMod := I1NNN(emitter.currentAddress+10)
	err = emitter.saveOpcode(skipMod)

	err = emitter.saveOpcode(I6XKK(1, 255)) //v1 = 255. We can use it as a helper because both operands are simples in the context of %


	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I6XKK(0xf, 0))	 // Vf = 0


	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I8XY5(0, 2))  // V0 = V0-V2


	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode( I4XKK(0, 0))	 //if v0 != 0 we skip the next opcode


	if err != nil{
		return 0, err
	}

	//so if v0 =0, we need stop dividing and we jump to the end
	jumpToEnd := I1NNN(emitter.currentAddress+5)
	err = emitter.saveOpcode(jumpToEnd)

	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I3XKK(0xf, 0))	//if not we ask if v0>v2


	if err != nil{
		return 0, err
	}
	//if v0>v2 we keep dividing in loop by jumping
	loop := I1NNN(emitter.currentAddress-5)
	err = emitter.saveOpcode(loop)

	if err != nil{
		return 0, err
	}

	//if not we jump the previous opcode and we find the rest by subtracting 255 (saved in v1) and v0.
	//That give us the rest
	err = emitter.saveOpcode(I8XY5(1, 0))//  = V1-V0


	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I8XY5(0, 1))  //  = V0 = V1 to save the rest in v0


	if err != nil{
		return 0, err
	}

	return 1, nil
}


//division translates a / to opcodes and write it in emitter.machineCode, return the size of the datatype of the result and an error
func (emitter *Emitter) division(functionCtx *FunctionCtx) (int, error) {
	_, err := emitter.solveOperands(functionCtx)
	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I6XKK(1, 0)) //v1 = 0. We can use it to store the result because both operands are simples in the context of /


	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I4XKK(0, 0)) //if v0 != 0 we skip the next opcode


	if err != nil{
		return 0, err
	}
	//if v0 =0, the result is 0 and we skip the division
	skipDivision := I1NNN(emitter.currentAddress+12)
	err = emitter.saveOpcode(skipDivision)

	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I6XKK(0xf, 0))// Vf = 0

	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I8XY5(0, 2))	 // V0 = V0-V2


	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I4XKK(0, 0))	 //if v0 != 0 we skip the next opcode


	if err != nil{
		return 0, err
	}


	err = emitter.saveOpcode(I7XKK(1,1))	//if v0 = 0 we do v1 = v1 + 1, to operate before jumping


	if err != nil{
		return 0, err
	}
	err = emitter.saveOpcode(I4XKK(0, 0))  //if v0 != 0 we skip the next opcode


	if err != nil{
		return 0, err
	}

	//if v0 =0, the rest of division is also 0 and we jump to the end of the operation
	jumpToEnd := I1NNN(emitter.currentAddress+5)
	err = emitter.saveOpcode(jumpToEnd)

	if err != nil{
		return 0, err
	}
	//if not we ask if v0>v2, and if v0 > v2 we skip the next opcode

	err = emitter.saveOpcode(I3XKK(0xf, 0))

	if err != nil{
		return 0, err
	}
	//if v0<v2 we jump to to the end of the division, if not we keep dividing
	jumpToEnd = I1NNN(emitter.currentAddress+2)
	err = emitter.saveOpcode(jumpToEnd)

	if err != nil{
		return 0, err
	}

	err = emitter.saveOpcode(I7XKK(1,1))

	if err != nil{
		return 0, err
	}

	loop := I1NNN(emitter.currentAddress-9)
	err = emitter.saveOpcode(loop)

	if err != nil{
		return 0, err
	}


	err = emitter.saveOpcode(I8XY5(0, 1))//  = V0 = V1 to save the result in v0


	if err != nil{
		return 0, err
	}

	return 1, nil
}

//solveOperands save the operands of a operation in registers. It return the indexes of registers in which each operand
//was stored and an error if needed
func (emitter *Emitter) solveOperands(functionCtx *FunctionCtx) (*ResultRegIndex, *ResultRegIndex, error){
	leftOperand := emitter.ctxNode.Children[0]
	rightOperand := emitter.ctxNode.Children[1]
	backup := emitter.ctxNode
	emitter.ctxNode = rightOperand
	var err error
	rightOperandRegIndex, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return nil, nil, err
	}

	emitter.ctxNode = leftOperand
	leftOperandRegIndex, err := emitter.translateOperation[emitter.ctxNode.Value.Type](functionCtx)
	if err != nil{
		return  nil, nil, err
	}

	emitter.ctxNode = backup
	return leftOperandRegIndex, rightOperandRegIndex, nil

}

//index save in registers the value of a dereference.
//Returns the indexes of registers in which it was stored the dereference and an error if needed
func (emitter *Emitter)index(functionCtx *FunctionCtx) (*ResultRegIndex, error){
	return emitter.saveDereferenceInRegisters(functionCtx)
}

//asterisk multiply registers or save a dereference in registers, depending on the context
//it returns the indexes of registers in which the result of the operation was stored and an error if needed
func (emitter *Emitter)asterisk(functionCtx *FunctionCtx)(*ResultRegIndex, error){
	if len(emitter.ctxNode.Children) == 1{
		return emitter.saveDereferenceInRegisters(functionCtx)
	}else{
		return emitter.multiplication(functionCtx)
	}

}

//saveDereferenceInRegisters save in registers  the value of a dereference.
//Returns the indexes of the registers in which it was stored and an error
func (emitter *Emitter) saveDereferenceInRegisters(functionCtx *FunctionCtx) (*ResultRegIndex, error) {
	size, err := emitter.saveDereferenceAddressInI(functionCtx)
	if err != nil {
		return nil, err
	}
	err = emitter.saveOpcode(IFX65(byte(size-1)))
	if err != nil{
		return nil, err
	}

	return emitter.allocAndCopyPaste(functionCtx, size, 0, 1)
}
//ident save registers the value of a reference.
//Returns the indexes of registers that use to save its values and an error if needed
func (emitter *Emitter)ident(functionCtx *FunctionCtx) (*ResultRegIndex, error){
	ident := emitter.ctxNode.Value.Literal
	var size int
	var err error

	_, isGlobalReference := emitter.globalVariables[ident]
	if isGlobalReference{
		size, err = emitter.saveGlobalReferenceAddressInI(0,1)
		if err != nil {
		return nil, err
		}
	}else{
		size, err = emitter.saveStackReferenceAddressInI(0,functionCtx)
		if err != nil{
			return nil, err
		}
	}
	err = emitter.saveOpcode(IFX65(byte(size-1)))
	if err != nil{
		return nil, err
	}
	regIndex, err := emitter.allocAndCopyPaste(functionCtx, size, 0, 1) //we save the value of the reference in available registers
	if err != nil {
		return regIndex, err
	}

	return regIndex, nil

}

//allocAndCopyPaste check the size of a variable (saved in vx and vy) and store it in registers.
//It return the index of this registers and an error if needed
func (emitter *Emitter) allocAndCopyPaste(functionCtx *FunctionCtx, size int, x byte, y byte) (*ResultRegIndex, error) {
	var regIndex *ResultRegIndex
	var ok bool
	switch size {
	case 1:
		regIndex, ok = functionCtx.registerHandler.AllocSimple()
		err := emitter.saveOpcode(I8XY0(regIndex.lowBitsIndex, x))
		if err != nil {
			return nil, err
		}

	case 2:
		regIndex, ok = functionCtx.registerHandler.AllocPointer()
		err := emitter.saveOpcode(I8XY0(regIndex.highBitsIndex, x))
		if err != nil {
			return nil, err
		}
		err = emitter.saveOpcode(I8XY0(regIndex.lowBitsIndex, y))
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}
	if !ok {
		line := emitter.ctxNode.Value.Line
		err := errors.New(errorhandler.TooManyRegisters(line))
		return nil, err
	}

	return regIndex, nil
}

//address save the address of its children in two registers, return the indexes of registers in which
//it stores it and an error if needed
func (emitter *Emitter)address(functionCtx *FunctionCtx)(*ResultRegIndex, error){
	emitter.ctxNode = emitter.ctxNode.Children[0]
	regIndex, ok := functionCtx.registerHandler.AllocPointer()
	if !ok{
		line := emitter.ctxNode.Value.Line
		err := errors.New(errorhandler.TooManyRegisters(line))
		return nil, err
	}

	//we save the address in I
	if emitter.ctxNode.Value.Type == token.IDENT{
		ident := emitter.ctxNode.Value.Literal
		_, isGlobalReference := emitter.globalVariables[ident]
		if isGlobalReference{
			_, err := emitter.saveGlobalReferenceAddressInI(0,1)
			if err != nil{
				return nil, err
			}
		}else{
			_, err := emitter.saveStackReferenceAddressInI(0,functionCtx)
			if err != nil{
				return nil, err
			}
		}
	}else{
		_,err := emitter.saveDereferenceAddressInI(functionCtx)
		if err != nil{
			return nil, err
		}
	}
	//then we save i in the registers
	err := emitter.saveOpcode(I9XY2(regIndex.highBitsIndex,regIndex.lowBitsIndex))
	if err != nil{
		return nil, err
	}
	return 	regIndex, nil
}

//saveGlobalReferenceAddressInI saves the address of a global variable in I using the register x and y
//Returns the size of the reference it points to and an error
func (emitter *Emitter)saveGlobalReferenceAddressInI(x byte, y byte) (int, error){
	ident := emitter.ctxNode.Value.Literal
	address := emitter.globalVariables[ident]
	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType)

	err :=emitter.saveOpcode(I6XKK(x, byte(address<<8)))
	if err != nil{
		return 0,err
	}

	err =emitter.saveOpcode( I6XKK(y, byte(address)))
	if err != nil{
		return 0,err
	}

	err =emitter.saveOpcode(I9XY1(x,y))
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
	_, isInStack := functionCtx.stack.References[leafIdent]
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
			err = emitter.saveFX1ESafely(0, index*symboltable.GetSize(datatype))
			datatype = datatype.(symboltable.Array).Of
			emitter.ctxNode = emitter.ctxNode.Children[1]

		//if we are analyzing a *, then its value is the address  of the next referenced element, si we set I = value.
		case token.ASTERISK:
			//we set V0 and V1 = value saved from I in memory
			err :=emitter.saveOpcode(IFX65(1))
			if err != nil{
				return 0,err
			}
			//we set I=value founded previously in I

			err =emitter.saveOpcode(I9XY1(0,1))
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
	reference, ok := functionCtx.stack.References[ident]
	if !ok{
		return 0, errors.New(errorhandler.UnexpectedCompilerError())
	}
	size := symboltable.GetSize(emitter.scope.Symbols[ident].DataType)

	//we set I = address position 0 of stack

	err := emitter.saveOpcode(I9XY1(RegisterStackAddress1, RegisterStackAddress2))
	if err != nil{
		return 0, err
	}

	return size, emitter.saveFX1ESafely(x, reference.positionInStack)

}

//saveFX1ESafely set an int to vx and then set I = I + vx, if the int is greater than 255 we add vx in a loop
func (emitter *Emitter) saveFX1ESafely(x byte, vx int) error{

	for  vx>255{


		err := emitter.saveOpcode(I6XKK(x, 255))
		if err != nil{
			return err
		}

		err =emitter.saveOpcode( IFX1E(x))
		if err != nil{
			return err
		}
		vx = vx - 255
	}
	if vx > 0{


		err :=emitter.saveOpcode(I6XKK(x, byte(vx)))
		if err != nil{
			return err
		}

		err =emitter.saveOpcode( IFX1E(x))
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