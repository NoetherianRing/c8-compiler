package emitter

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
)

type Emitter struct{
//	memoryHandler   *MemoryHandler
//	registerHandler *RegistersHandler
	lastAvailableAddress uint16
	globalVariables		map[string]uint16
	scope      			 *symboltable.Scope
	mainScope      		*symboltable.Scope
	ctxNode     		*ast.Node
	machineCode 		[MEMORY]byte
	translate 		    map[string]func()error
	functions   		map[string]uint16
}

func NewEmitter(tree *ast.SyntaxTree, scope *symboltable.Scope)*Emitter{
	emitter := new(Emitter)
//	emitter.memoryHandler = NewMemoryHandler()
//	emitter.registerHandler = NewRegisterHandler()
	emitter.globalVariables = make(map[string]uint16)
	emitter.scope = scope
	emitter.mainScope = scope
	emitter.ctxNode = tree.Head
	emitter.translate = make(map[string]func()error)
	emitter.lastAvailableAddress = GlobalSectionStart
	//	emitter.translate[token.EOF] = emitter.eof
	return emitter
}


func (emitter *Emitter) Start() ([MEMORY]byte, error){
	var globalDeclSection []byte
	emitter.ctxNode = emitter.ctxNode.Children[0] //The tree start with a EOF node, so we move to the next one
	err := emitter.primitivesFn()
	if err != nil{
		return emitter.machineCode, err
	}
	//first we save all the global variables into memory
	block := emitter.ctxNode
	for _, child := range block.Children{
		if child.Value.Type == token.LET{
			emitter.ctxNode = child
			globalDeclSection, err = emitter.globalDecl()
		 	// err = emitter.globalDecl()
			if err != nil{
				return emitter.machineCode, err
			}
		}
	}
	emitter.ctxNode = block
	if len(globalDeclSection) > MEMORY - GlobalSectionStart{
			return emitter.machineCode, errors.New(errorhandler.NotEnoughMemory())
	}
	for i, variable := range globalDeclSection{
		emitter.machineCode[GlobalSectionStart+i] = variable
		emitter.lastAvailableAddress++
	}

	//we save into memory the primitive functions
	err = emitter.primitivesFn()
	if err != nil{
		return emitter.machineCode, err
	}
	//we save into memory the functions
	for _, child := range block.Children{
		if child.Value.Type == token.FUNCTION{

			err = emitter.fn()
			if err != nil{
				return emitter.machineCode, err
			}
		}
	}
	emitter.ctxNode = block

	//The stack section will start in the last available address, which is saved in the v2 and v3 registers
	v2 := byte(emitter.lastAvailableAddress & 0xFF00 >> 8)
	v3 := byte(emitter.lastAvailableAddress & 0x00FF)

	saveV2 := I6XKK(FirstRegisterForStackSection, v2)
	emitter.machineCode[0] = saveV2[0]
	emitter.machineCode[1] = saveV2[1]

	saveV3 := I6XKK(SecondRegisterForStackSection, v3)
	emitter.machineCode[2] = saveV3[0]
	emitter.machineCode[3] = saveV3[1]

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

func (emitter *Emitter)primitivesFn()error{
	return nil

}
func (emitter *Emitter)fn()error{
	return nil
}

func (emitter *Emitter)globalDecl()([]byte, error){
	let := emitter.ctxNode
	ident := let.Children[0].Value.Literal
	symbol, ok := emitter.scope.Symbols[ident]
	if !ok{
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}

	machineCode, err := emitter.declaration(symbol.Identifier, symbol.DataType, emitter.globalVariables, emitter.lastAvailableAddress)
	if err != nil{
		return nil, err
	}

	return machineCode, nil
}


func (emitter *Emitter) declaration(ident string, datatype interface{}, variables map[string]uint16, address uint16)([]byte, error){
	_, exist := variables[ident]
	if !exist{
		variables[ident] = address
	}

	switch datatype.(type) {
	case symboltable.Simple:
		 machineCode := make([]byte,1)
		 return machineCode, nil
	case symboltable.Pointer:
		return emitter.declarePointer(datatype.(symboltable.Pointer), address)
	case symboltable.Array:
		 return emitter.declareArray(datatype.(symboltable.Array), address)
	default:
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}
}

func (emitter *Emitter) complexDeclaration(datatype interface{}, address uint16)([]byte, error) {
	switch datatype.(type) {
	case symboltable.Pointer:
		return emitter.declarePointer(datatype.(symboltable.Pointer), address)
	case symboltable.Array:
		return emitter.declareArray(datatype.(symboltable.Array), address)
	default:
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}
}

//declarePointer
func (emitter *Emitter)declarePointer(pointer symboltable.Pointer, address uint16)([]byte, error) {
	machineCode := make([]byte,0)
	pointsToAddress := address + 2
	machineCode = append(machineCode, byte(pointsToAddress>>8))
	machineCode = append(machineCode, byte(pointsToAddress))
/*	emitter.saveValueInAddress(machineCode, address, byte(pointsToAddress>>8)) //we save in the given address
	// the most significant 8 bits of the address it points to
	emitter.saveValueInAddress(machineCode, address+1, byte(pointsToAddress)) //we save in the the next address
	// the least significant 8 bits of the address it points to
*/
	switch pointer.PointsTo.(type) {
	case symboltable.Simple:
		return append(machineCode,0), nil
	default:
		nextDecl, err := emitter.complexDeclaration(pointer.PointsTo,  pointsToAddress)
		if err != nil{
			return nil, err
		}
		return append(machineCode,nextDecl...), nil

	}

}
//TODO: Revisar
func (emitter *Emitter)declareArray(array symboltable.Array, address uint16)([]byte,error)  {
	switch array.Of.(type) {
	case symboltable.Simple:
		machineCode := make([]byte, array.Length)
		return machineCode, nil
	default:

		machineCode := make([]byte,0)

		for i:=0; i<array.Length;i++{
			//pointsToAddress := address + 2
			//machineCode = append(machineCode, byte(pointsToAddress>>8))
			//machineCode = append(machineCode, byte(pointsToAddress))
			//nextDecl, err := emitter.complexDeclaration(array.Of,  pointsToAddress)
			nextDecl, err := emitter.complexDeclaration(array.Of,  address)
			if err != nil{
				return nil, err
			}
			address = address + uint16(len(nextDecl))
			machineCode = append(machineCode,nextDecl...)
		}
		return machineCode, nil
	}
}

func (emitter *Emitter)saveValueInAddress(writeIn []byte, address uint16, value byte) {
}
func (emitter *Emitter)createStack()[]byte{
	return nil

}

func (emitter *Emitter)deleteStack()[]byte{
	return nil

}

