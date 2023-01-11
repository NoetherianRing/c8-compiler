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
	analyzer.currentScope = symboltable.CreateMainScope()
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

func (analyzer *SemanticAnalyzer) Start() error{
	ok := analyzer.savePrimitiveFunctions()
	if !ok{
		panic(errorhandler.UnexpectedCompilerError())
	}
	next := analyzer.ctxNode.Children[0].Value.Type
	return analyzer.validate[next]()
}

func(analyzer *SemanticAnalyzer) block()error{
	analyzer.currentScope.AddSubScope()
	lastAdded := len(analyzer.currentScope.SubScopes)-1
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded]
	backup := analyzer.ctxNode
	for _, child := range analyzer.ctxNode.Children{
		analyzer.ctxNode = child
		next := child.Value.Type
		err := analyzer.validate[next]()
		if err != nil{
			return err
		}
	}
	analyzer.ctxNode = backup
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded].Parent
	return nil
}

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
func(analyzer *SemanticAnalyzer) assign()error{
	leftTree := analyzer.ctxNode.Children[0]
	analyzer.updateDataTypeFactoryCtx(leftTree)
	leftDataType, err :=analyzer.datatypeFactory.GetDataType()
	if err != nil{
		return err
	}
	rightTree := analyzer.ctxNode.Children[0]
	analyzer.updateDataTypeFactoryCtx(rightTree)
	rightDataType, err2 :=analyzer.datatypeFactory.GetDataType()

	if err2 != nil{
		return err
	}

	if !symboltable.Compare(leftDataType, rightDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(leftDataType), token.EQ, symboltable.Fmt(rightDataType)))
		return err
	}
	return nil
}

func(analyzer *SemanticAnalyzer) fn()error{
	backup := analyzer.ctxNode
	analyzer.ctxNode = analyzer.ctxNode.Children[0]
	args, err := analyzer.handleParams()
	if err != nil{
		return err
	}
	analyzer.ctxNode = backup

	name := analyzer.ctxNode.Children[1].Value.Literal
	analyzer.updateDataTypeFactoryCtx(analyzer.ctxNode.Children[2])

	expectedReturnDataType, err2 := analyzer.datatypeFactory.GetDataType()
	if err2 != nil{
		return err2
	}
	analyzer.ctxNode = analyzer.ctxNode.Children[3]

	actualReturnDataType, err3 := analyzer.funcBlock()
	if err3 != nil{
		return err2
	}

	if !symboltable.Compare(expectedReturnDataType, actualReturnDataType){
		line := analyzer.ctxNode.Value.Line
		err = errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(expectedReturnDataType), token.EQ, symboltable.Fmt(actualReturnDataType)))
		return err
	}
	analyzer.ctxNode = backup
	function := symboltable.NewFunction(expectedReturnDataType, args)
	analyzer.currentScope.AddSymbol(name, function)
	return nil

}

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

func(analyzer *SemanticAnalyzer) _if()error{
	return analyzer.validateConditionAndBlock()
}

func(analyzer *SemanticAnalyzer) _else()error{
	err := analyzer.validateConditionAndBlock()
	if err != nil{
		return err
	}
	elseBlock :=  analyzer.ctxNode.Children[2].Value.Type
	err = analyzer.validate[elseBlock]()
	return err
}

func(analyzer *SemanticAnalyzer) _while()error{
	return analyzer.validateConditionAndBlock()
}

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

func(analyzer *SemanticAnalyzer) savePrimitiveFunctions() bool{
	return analyzer.saveDraw()  &&
		analyzer.saveSetDT() && analyzer.saveGetDT() &&
		analyzer.saveSetST() && analyzer.saveGetST()
}

func(analyzer *SemanticAnalyzer) saveDraw() bool{
	byteType := symboltable.NewByte()
	paramType := make([]interface{},4)
	paramType[0] = byteType //x
	paramType[1] = byteType//y
	paramType[2] = symboltable.NewPointer(byteType) //sprite address
	paramType[3] = byteType //length
	returnType := symboltable.NewBool() //collision
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("DRAW", functionType)
}

func(analyzer *SemanticAnalyzer) saveSetDT() bool{
	paramType := make([]interface{},1)
	paramType[0] = symboltable.NewByte()
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("SET_DT", functionType)
}

func(analyzer *SemanticAnalyzer) saveGetDT() bool{
	returnType := symboltable.NewByte()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("GET_DT", functionType)
}

func(analyzer *SemanticAnalyzer) saveSetST() bool{
	paramType := make([]interface{},1)
	paramType[0] = symboltable.NewByte()
	returnType := symboltable.NewVoid()
	functionType := symboltable.NewFunction(returnType, paramType)
	return analyzer.currentScope.AddSymbol("SET_ST", functionType)
}

func(analyzer *SemanticAnalyzer) saveGetST() bool{
	returnType := symboltable.NewByte()
	functionType := symboltable.NewFunction(returnType, nil)
	return analyzer.currentScope.AddSymbol("GET_ST", functionType)
}

func(analyzer *SemanticAnalyzer) updateDataTypeFactoryCtx(toAnalyze *ast.Node) {
	analyzer.datatypeFactory.SetCxtNode(toAnalyze)
	analyzer.datatypeFactory.SetScope(analyzer.currentScope)
}

func(analyzer *SemanticAnalyzer) funcBlock()(interface{}, error){
	analyzer.currentScope.AddSubScope()
	lastAdded := len(analyzer.currentScope.SubScopes)-1
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded]
	backup := analyzer.ctxNode
	for _, child := range analyzer.ctxNode.Children{
		if child.Value.Type != token.RETURN{
			analyzer.ctxNode = child
			next := analyzer.ctxNode.Value.Type
			err := analyzer.validate[next]()
			analyzer.ctxNode = backup
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
	analyzer.currentScope = analyzer.currentScope.SubScopes[lastAdded].Parent
	return symboltable.NewVoid(), nil
}

func (analyzer *SemanticAnalyzer) handleParams()([]interface{},error) {
	if len(analyzer.ctxNode.Children) == 0 {
		return nil, nil
	}
	args := make([]interface{},0)
	analyzer.ctxNode = analyzer.ctxNode.Children[0]
	for len(analyzer.ctxNode.Children) == 2{
		 param, err := analyzer.handleParam(0)
		 if err != nil{
		 	return nil, err
		 }
    	args = append(args, param)

    	if analyzer.ctxNode.Children[1].Value.Type != token.COMMA{
			param, err = analyzer.handleParam(1)

			if err != nil{
			  return nil, err
			}
			args = append(args, param)
		}else{
			analyzer.ctxNode = analyzer.ctxNode.Children[1]

		}

	}
	if len(analyzer.ctxNode.Children) == 1{
		param, err := analyzer.handleParam(0)
		if err != nil{
			return nil, err
		}

		args = append(args, param)
	}
	return args, nil
}

func (analyzer *SemanticAnalyzer) handleParam(i int) (interface{}, error) {
	backup := analyzer.ctxNode
	analyzer.ctxNode = analyzer.ctxNode.Children[i]
	err := analyzer.validate[analyzer.ctxNode.Value.Type]()
	if err != nil {
		return nil, err
	}
	analyzer.ctxNode = backup
	analyzer.updateDataTypeFactoryCtx(analyzer.ctxNode.Children[i])
	datatype, err2 := analyzer.datatypeFactory.GetDataType()
	return datatype, err2

}