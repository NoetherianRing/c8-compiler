package semanticAnalyzer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
)

type statementValidator map[token.Type]func()error

type SemanticAnalyzer struct{
	datatypeFactory *DataTypeFactory
	validate        statementValidator
	currentScope    *symboltable.Scope
	ctxNode         *ast.Node
}

func NewSemanticAnalyzer(tree *ast.SyntaxTree)*SemanticAnalyzer{
	analyzer := new(SemanticAnalyzer)
	analyzer.datatypeFactory = NewDataTypeFactory()
	analyzer.currentScope = symboltable.CreateGlobalScope()
	analyzer.ctxNode = tree.Head
	analyzer.validate = make(statementValidator)
	analyzer.validate[token.RBRACE] = analyzer.block
	analyzer.validate[token.LET] = analyzer.let
	analyzer.validate[token.EQ] = analyzer.assign
	analyzer.validate[token.FUNCTION] = analyzer.fn
	analyzer.validate[token.RPAREN] = analyzer.call
	analyzer.validate[token.IF] = analyzer._if
	analyzer.validate[token.ELSE] = analyzer._else
	analyzer.validate[token.WHILE] = analyzer._while
	return analyzer
}

//Start save in the symbol table the primitive functions, and validates the semantic of global declarations
//It also checks the declaration of a main function
func (analyzer *SemanticAnalyzer) Start() (*symboltable.Scope, error){
	ok := analyzer.savePrimitiveFunctions()
	if !ok{
		panic(errorhandler.UnexpectedCompilerError())
	}
	globalScope := analyzer.currentScope
	block := analyzer.ctxNode.Children[0].Children[0] //The tree start with a "" and a EOF node, so we move

	for _, declaration := range block.Children{
		analyzer.ctxNode = declaration
		next := declaration.Value.Type
		if next != token.FUNCTION && next != token.LET{
			line := analyzer.ctxNode.Value.Line
			return globalScope, errors.New(errorhandler.GlobalScopeOnlyAllowsDeclarations(line))
		}
		err := analyzer.validate[next]()
		if err != nil{
			return globalScope, err
		}
	}
	_, existMain := globalScope.Symbols[token.MAIN]

	if !existMain{
		return globalScope, errors.New(errorhandler.MainFunctionNeeded())
	}

	return globalScope, nil
}

//block creates a new sub scope and validates the semantic of all the statements within the block
func(analyzer *SemanticAnalyzer) block()error{
	analyzer.currentScope.AddSubScope()
	lastAdded := len(analyzer.currentScope.SubScopes)-1
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded]
	block := analyzer.ctxNode
	for _, child := range block.Children{
		analyzer.ctxNode = child
		next := analyzer.ctxNode.Value.Type
		//functions only can be declared in the global scope
		if next == token.FUNCTION{
			line := analyzer.ctxNode.Value.Line
			return errors.New(errorhandler.FunctionOutsideGlobalScope(line))
		}
		err := analyzer.validate[next]()
		if err != nil{
			return err
		}
	}
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded].Parent
	return nil
}

//let validates the semantic of a declaration statements, checks that the name of the declaration is not already in use,
//and if its not, save the new variable in the symbol table of the current scope
func(analyzer *SemanticAnalyzer) let()error{
	name := analyzer.ctxNode.Children[0].Value.Literal
	datatypeTree := analyzer.ctxNode.Children[1]
	analyzer.updateDataTypeFactoryCtx(datatypeTree)
	datatype, err := analyzer.datatypeFactory.GetDataType()
	if err != nil{
		return err
	}
	ok := analyzer.currentScope.AddSymbol(name, datatype)
	if !ok{
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.NameAlreadyInUse(line, name))
		return err
	}

	return nil
}

//assign validates the semantic of assignation statements
func(analyzer *SemanticAnalyzer) assign()error{
	leftTree := analyzer.ctxNode.Children[0]
	analyzer.updateDataTypeFactoryCtx(leftTree)
	leftDataType, err :=analyzer.datatypeFactory.GetDataType()

	if err != nil{
		return err
	}
	if symboltable.IsAnArray(leftDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.InvalidAssignation(line, symboltable.Fmt(leftDataType)))
		return err

	}

	rightTree := analyzer.ctxNode.Children[0]
	analyzer.updateDataTypeFactoryCtx(rightTree)
	rightDataType, err :=analyzer.datatypeFactory.GetDataType()

	if err != nil{
		return err
	}

	if !symboltable.Compare(leftDataType, rightDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(leftDataType), token.EQ, symboltable.Fmt(rightDataType)))
		return err
	}
	return nil
}

//fn validates the semantic of the declaration of a function, then checks that the name of the declaration is not already in use,
//and if its not, save the new variable in the symbol table of the current scope
func(analyzer *SemanticAnalyzer) fn()error{
	const IDENT =0
	const PARAMS =1
	const DATATYPERETURN =2
	const BLOCK =3
	backupNode := analyzer.ctxNode
	backupScope := analyzer.currentScope
	analyzer.ctxNode = analyzer.ctxNode.Children[PARAMS]
	args, err := analyzer.handleParams()
	if err != nil{
		return err
	}
	analyzer.ctxNode = backupNode

	name := analyzer.ctxNode.Children[IDENT].Value.Literal
	analyzer.updateDataTypeFactoryCtx(analyzer.ctxNode.Children[DATATYPERETURN])

	expectedReturnDataType, err := analyzer.datatypeFactory.GetDataType()
	if err != nil{
		return err
	}

	if !symboltable.Compare(expectedReturnDataType, symboltable.NewVoid()) &&
		!symboltable.Compare(expectedReturnDataType, symboltable.NewBool()) &&
		!symboltable.Compare(expectedReturnDataType, symboltable.NewByte()){

		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.InvalidReturnType(line, symboltable.Fmt(expectedReturnDataType)))
		return err

	}
	analyzer.ctxNode = analyzer.ctxNode.Children[BLOCK]

	actualReturnDataType, err := analyzer.funcBlock()
	if err != nil{
		return err
	}

	if !symboltable.Compare(expectedReturnDataType, actualReturnDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(expectedReturnDataType), token.EQ, symboltable.Fmt(actualReturnDataType)))
		return err
	}
	analyzer.ctxNode = backupNode
	analyzer.currentScope = backupScope

	function := symboltable.NewFunction(expectedReturnDataType, args)
	ok := analyzer.currentScope.AddSymbol(name, function)
	if !ok{
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.NameAlreadyInUse(line, name))
	}

	return nil

}

//call validates the semantic of functions calls
func(analyzer *SemanticAnalyzer) call()error{
	toAnalyze := analyzer.ctxNode
	analyzer.updateDataTypeFactoryCtx(toAnalyze)
	funcDataType, err := analyzer.datatypeFactory.GetDataType()
	if err != nil{
		return err
	}
	if symboltable.NewVoid().Compare(funcDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.UnreachableCode(line))
		return err

	}
	return nil
}

//_if validates the semantic of if statements
func(analyzer *SemanticAnalyzer) _if()error{
	return analyzer.validateConditionAndBlock()
}

//_else validates the semantic of if/else statements
func(analyzer *SemanticAnalyzer) _else()error{
	err := analyzer.validateConditionAndBlock()
	if err != nil{
		return err
	}
	elseBlock :=  analyzer.ctxNode.Children[2].Value.Type
	err = analyzer.validate[elseBlock]()
	return err
}

//_while validates the semantic of while statements
func(analyzer *SemanticAnalyzer) _while()error{
	return analyzer.validateConditionAndBlock()
}

//validateConditionAndBlock validates that the condition of a statement suck as if, if/else and while
//is a boolean expression, and then executes the block of the statement
func (analyzer *SemanticAnalyzer) validateConditionAndBlock() error {
	boolDatatype := symboltable.NewBool()
	condition := analyzer.ctxNode.Children[0]

	analyzer.updateDataTypeFactoryCtx(condition)
	datatypeCondition, err := analyzer.datatypeFactory.GetDataType()
	if err != nil {
		return err
	}
	if boolDatatype.Compare(datatypeCondition) {
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.UnexpectedDataType(line, symboltable.Fmt(boolDatatype), symboltable.Fmt(datatypeCondition)))
		return err
	}
	block := analyzer.ctxNode.Children[1].Value.Type
	err = analyzer.validate[block]()
	return err
}

//savePrimitiveFunctions save into the symbol table the primitive functions of the language
func(analyzer *SemanticAnalyzer) savePrimitiveFunctions() bool{
	return analyzer.saveDraw()  && analyzer.saveClean() &&
		analyzer.saveSetDT() && analyzer.saveGetDT() &&
		analyzer.saveSetST() && analyzer.saveWaitKey() &&
		analyzer.saveDrawFont() && analyzer.saveIsKeyPressed() &&
		analyzer.saveRandom()
}

//saveDraw save into the symbol table a function named Draw that represents the chip-8 opcode DXYN
func(analyzer *SemanticAnalyzer) saveDraw() bool{
	byteType := symboltable.NewByte()
	paramType := make([]interface{},4)
	paramType[0] = byteType //x
	paramType[1] = byteType//y
	paramType[2] = byteType //length
	paramType[3] = symboltable.NewPointer(byteType) //sprite address
	returnType := symboltable.NewBool() //collision
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("Draw", functionType)
}

//saveDrawFont save into the symbol table a function named DrawFont that represents the chip-8 opcode DXYN with I = font
func(analyzer *SemanticAnalyzer) saveDrawFont() bool{
	byteType := symboltable.NewByte()
	paramType := make([]interface{},4)
	paramType[0] = byteType //x
	paramType[1] = byteType//y
	paramType[2] = byteType //font
	returnType := symboltable.NewBool() //collision
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("DrawFont", functionType)
}

//saveClean save into the symbol table a function named Clean that represents the chip-8 opcode I00E0
func(analyzer *SemanticAnalyzer) saveClean() bool{
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("Clean", functionType)
}


//saveSetDT save into the symbol table a function named SetDT that represents the chip-8 opcode FX15
func(analyzer *SemanticAnalyzer) saveSetDT() bool{
	paramType := make([]interface{},1)
	paramType[0] = symboltable.NewByte()
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("SetDT", functionType)
}

//saveGetDT save into the symbol table a function named GetDT that represents the chip-8 opcode FX07
func(analyzer *SemanticAnalyzer) saveGetDT() bool{
	returnType := symboltable.NewByte()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("GetDT", functionType)
}

//saveSetST save into the symbol table a function named SetST that represents the chip-8 opcode FX18
func(analyzer *SemanticAnalyzer) saveSetST() bool{
	paramType := make([]interface{},1)
	paramType[0] = symboltable.NewByte()
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("SetST", functionType)
}

//saveRandom save into the symbol table a function named Random that represents the chip-8 opcode CXKK
func(analyzer *SemanticAnalyzer) saveRandom() bool{
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("Random", functionType)
}

//saveWaitKey save into the symbol table a function named WaitKey that represents the chip-8 opcode FX0A
func(analyzer *SemanticAnalyzer) saveWaitKey() bool{
	returnType := symboltable.NewByte()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("WaitKey", functionType)
}
//saveIsKeyPressed save into the symbol table a function named IsKeyPressed that returns true if the key was pressed
func(analyzer *SemanticAnalyzer) saveIsKeyPressed() bool{
	paramType := make([]interface{},1)
	paramType[0] = symboltable.NewByte()
	returnType := symboltable.NewBool()
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("IsKeyPressed", functionType)
}

//updateDataTypeFactoryCtx updates the context of datatypeFactory
func(analyzer *SemanticAnalyzer) updateDataTypeFactoryCtx(toAnalyze *ast.Node) {
	analyzer.datatypeFactory.SetCxtNode(toAnalyze)
	analyzer.datatypeFactory.SetScope(analyzer.currentScope)
}

//funcBlock adds a new sub scope and validates the semantic of all the statements within the block,
//then returns the data type of the return statement
func(analyzer *SemanticAnalyzer) funcBlock()(interface{}, error){
	analyzer.currentScope.AddSubScope()
	lastAdded := len(analyzer.currentScope.SubScopes)-1
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded]
	block := analyzer.ctxNode
	for _, child := range block.Children{
		if child.Value.Type != token.RETURN{
			analyzer.ctxNode = child
			next := analyzer.ctxNode.Value.Type
			err := analyzer.validate[next]()
			if err != nil{
				return nil, err
			}
		}else{
			if len(child.Children) != 0{
				toAnalyze := child.Children[0]
				analyzer.updateDataTypeFactoryCtx(toAnalyze)
				return analyzer.datatypeFactory.GetDataType()

			}
		}
	}
	return symboltable.NewVoid(), nil
}

//handleParams validates the semantic of all the params of a function and save them in the symbol table of a new scope
//then returns an array with all the data types of the params
func (analyzer *SemanticAnalyzer) handleParams()([]interface{},error) {
	if len(analyzer.ctxNode.Children) == 0 { //ctxNode = )
		return nil, nil
	}

	args := make([]interface{},0)

	analyzer.ctxNode = analyzer.ctxNode.Children[0] //ctxNode = comma or let
	analyzer.currentScope.AddSubScope()
	lastAdded := len(analyzer.currentScope.SubScopes)-1
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded]
	totalSize := 0
	if analyzer.ctxNode.Value.Type == token.COMMA{
		for analyzer.ctxNode.Value.Type == token.COMMA{
			backup := analyzer.ctxNode
			analyzer.ctxNode = analyzer.ctxNode.Children[0] //ctxNode = let
			param, err := analyzer.handleParam()
			if err != nil{
				return nil, err
			}

			totalSize += symboltable.GetSize(param)
			args = append(args, param)
			analyzer.ctxNode = backup
			analyzer.ctxNode = analyzer.ctxNode.Children[1]

			if analyzer.ctxNode.Value.Type != token.COMMA{
				param, err = analyzer.handleParam()
				totalSize += symboltable.GetSize(param)

				if err != nil{
					return nil, err
				}

				args = append(args, param)
			}

		}
	}else{
		param, err := analyzer.handleParam()
		totalSize += symboltable.GetSize(param)

		if err != nil{
			return nil, err
		}

		args = append(args, param)
	}

	if totalSize > LimitParamSize{
		line := analyzer.ctxNode.Value.Line
		err := errors.New(errorhandler.ToManyParams(line))
		return nil, err
	}

	return args, nil
}

//handleParams validates the semantic of each single the param of a function and save it in the symbol table of a new scope
//then returns the data types of the param
func (analyzer *SemanticAnalyzer) handleParam() (interface{}, error) {
	backup := analyzer.ctxNode
	const IDENT = 0
	const DATATYPE = 1
	name := analyzer.ctxNode.Children[IDENT].Value.Literal
	datatypeTree := analyzer.ctxNode.Children[DATATYPE]
	analyzer.updateDataTypeFactoryCtx(datatypeTree)
	datatype, err := analyzer.datatypeFactory.GetDataType()
	if err != nil{
		return nil, err
	}
	ok := analyzer.currentScope.AddSymbol(name, datatype)
	if !ok{
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.NameAlreadyInUse(line, name))
		return nil, err
	}

	switch datatype.(type){
	case symboltable.Array:
		line := analyzer.ctxNode.Value.Line
		err := errors.New(errorhandler.InvalidParamType(line, symboltable.Fmt(datatype)))
		return nil, err
	}
	analyzer.ctxNode = backup

	return datatype, err

}