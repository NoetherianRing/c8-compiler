package semanticAnalyzer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
	"reflect"
	"strconv"
)



type DataTypeFactory struct{
	scope *symboltable.Scope
	//tree *ast.SyntaxTree
	ctxNode      *ast.Node
	walkingAFunc bool
	//execute map[token.Type]func()(interface{}, error)
}

func NewDataTypeGetter() *DataTypeFactory {
	getter := new(DataTypeFactory)
	getter.walkingAFunc = false
	getter.ctxNode = nil
	getter.scope = nil
	return getter
}

func (getter *DataTypeFactory)SetCxtNode(node *ast.Node){
	getter.ctxNode = node
}

func (getter *DataTypeFactory)SetScope(scope *symboltable.Scope){
	getter.scope = scope
}


//GetDataType calls "redirect()"to obtain a function that returns the data type of the current node of the tree.
//It then executes that function and returns its result.
func (getter *DataTypeFactory) GetDataType()(interface{}, error){
	if getter.ctxNode == nil || getter.scope == nil{
		panic(errorhandler.UnexpectedCompilerError())
	}
	if getter.scope == nil {
		panic(errorhandler.UnexpectedCompilerError())
	}
	get := getter.redirect()
	getter.scope = nil
	getter.ctxNode = nil
	return get()
}

//redirect analyzes the token type of the current node of the tree and returns a function that analyzes the data type of the expression led by the node.
//When needed, it calls another redirect method for further analysis.
func (getter *DataTypeFactory) redirect() func()(interface{}, error){


	switch getter.ctxNode.Value.Type{
	case token.ASTERISK:
		return getter.redirectAsterisk()
	case token.RPAREN:
		return getter.redirectParentheses()
	case token.VOID:
		return getter.declarationSimple
	case  token.TYPEBYTE:
		return getter.declarationSimple
	case token.TYPEBOOL:
		return getter.declarationSimple
	case token.IDENT:
		return getter.reference
	default:
		panic(errorhandler.UnexpectedCompilerError())
	}

}

//redirectAsterisk return a function that analysis the data type of an expression led by an asterisk by checking the context.
func (getter *DataTypeFactory) redirectAsterisk() func()(interface{}, error) {

	switch len(getter.ctxNode.Children){
	case 1:
		return getter.numericExpression
	case 2:
		if getter.isADeclarationContext(){
			return getter.declaration
		}else{
			return getter.dereference
		}
	default: panic(errorhandler.UnexpectedCompilerError())
	}

}

//redirectParentheses return a function that analysis the data type of an expression led by a parentheses by checking the context.
func (getter *DataTypeFactory) redirectParentheses() func()(interface{}, error) {
	child := getter.ctxNode.Children[0].Value
	if child.Type == token.IDENT{
		return getter.functionCall
	}
	return getter.skipNodeByLeft
}

//skipNodeByLeft skip the analysis of the data type of the current node, analysing its left child instead
func (getter *DataTypeFactory) skipNodeByLeft() (interface{}, error) {
	getter.ctxNode = getter.ctxNode.Children[0]
	return getter.GetDataType()
}

//isADeclarationContext checks if the leaf after a sequence of nodes has a value of token type "typebool" or "typebyte"
//to identify if the program is in a context of variable declaration.
func (getter *DataTypeFactory) isADeclarationContext() bool {
	//Due to the grammar and the syntax tree, the leaf that would provide information about the context
	//is going to be found by walking the tree using the right child of each node.
	leaf := GetLeafByRight(getter.ctxNode)
	leafType := leaf.Value.Type

	if leafType == token.TYPEBOOL|| leafType == token.TYPEBYTE{
		return true
	}
	return false
}


//validateParamsDataType validates if the data type of the parameters of a function call match with each data type
//of expressions led by the ctxNode. It also validates if the numbers of parameters matches.
func (getter *DataTypeFactory) validateParamsDataType(args []interface{}) error{
	var err error
	err = nil
	for i := range args{
		if getter.ctxNode.Value.Type == token.COMMA{

			comma := getter.ctxNode
			getter.ctxNode = getter.ctxNode.Children[0]
			err = getter.validateParamDataType(args, i)
			if err != nil{
				return err
			}

			getter.ctxNode = comma.Children[1]


		}else{
			if len(args) - (i) != 1{
				line := getter.ctxNode.Value.Line
				err =errors.New(errorhandler.NumberOfParametersDoesntMatch(line, i+1, len(args)))
			}else{
				err = getter.validateParamDataType(args, i)
			}
		}

	}
	if getter.ctxNode.Value.Type == token.COMMA {
		numberOfTreeParams := CountNodesToLeafByRight(getter.ctxNode)
		line := getter.ctxNode.Value.Line
		err =errors.New(errorhandler.NumberOfParametersDoesntMatch(line, numberOfTreeParams, len(args)))

	}

	return err
}

//validateParamDataType  validates if the data type of the expression led by the current "ctxNode"
//matches the data type of the argument "i" of a function.
func (getter *DataTypeFactory) validateParamDataType(args []interface{}, i int) error{
	treeParam, err := getter.GetDataType()
	if err != nil{
		return nil
	}
	if symboltable.Compare(treeParam, args[i]){
		return nil
	}else{

		line := getter.ctxNode.Value.Line
		err =errors.New(errorhandler.DataTypesDontMatch(line, symboltable.Fmt(treeParam), token.EQ, symboltable.Fmt(args[i])))
		return err
	}

}


func (getter *DataTypeFactory) numericExpression() (interface{}, error){
}

//functionCall validates that the data types of the parameters of a function call
//match the ones of its definition and returns an error if needed.
//If not, it returns the data type of the return value of the function.
func (getter *DataTypeFactory) functionCall() (interface{}, error){
	getter.walkingAFunc = true

	backup := getter.ctxNode
	identifier := getter.ctxNode.Children[0]


	getter.ctxNode = identifier
	identifierDataType, err := getter.reference()
	getter.ctxNode = backup

	if err != nil{
		return nil, err
	}

	funcDataType, ok := identifierDataType.(symboltable.Function)
	if !ok{
		panic(errorhandler.UnexpectedCompilerError())
	}

	returnDataType := funcDataType.Return
	argsDataType := funcDataType.Args

	if argsDataType != nil {
		if len(getter.ctxNode.Children) != 2 {
			line := getter.ctxNode.Value.Line
			err = errors.New(errorhandler.NumberOfParametersDoesntMatch(line, len(argsDataType), 0))
			return nil, err
		}

		param := getter.ctxNode.Children[1]
		getter.ctxNode = param
		err = getter.validateParamsDataType(argsDataType)
		getter.ctxNode = backup
		if err != nil{
			return nil, err
		}
	}
	getter.walkingAFunc = false
	return returnDataType, nil

}
//reference checks if an identifier is stored in the symbol table and, if it is, returns its data type.
//It returns an error in case of not finding the reference, if we are expecting a function and the reference is not a function,
//or if we are not expecting a function and the reference is a function.
func (getter *DataTypeFactory) reference()(interface{}, error){
	literal := getter.ctxNode.Value.Literal
	ref, ok := getter.scope.Symbols[literal]
	if !ok{
		line :=  getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnresolvedReference(line, literal))
		return nil, err
	}else{
		if getter.walkingAFunc && !ref.IsFunction{
				line :=  getter.ctxNode.Value.Line
				err := errors.New(errorhandler.IdentifierIsNotFunction(line, literal))
				return nil, err
		}

		if !getter.walkingAFunc && ref.IsFunction{
				line :=  getter.ctxNode.Value.Line
				err := errors.New(errorhandler.IdentifierIsFunction(line, literal))
				return nil, err
		}
	}
	return ref.DataType, nil
}

func (getter *DataTypeFactory) dereference() (interface{}, error){

}

//declaration verifies that there is no declaration of a pointer to a function
//in the context and returns an error if needed. Otherwise, it returns the data type built by "declarationBuild()".
func (getter *DataTypeFactory) declaration() (interface{}, error){
	if len(getter.ctxNode.Children) != 0{
		leaf := GetLeafByRight(getter.ctxNode)
		leafType := leaf.Value.Type
		if leafType == token.VOID {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.PointerToVoid(line))
			return nil, err
		}
	}
	return getter.declarationBuild()
}

//declarationBuild builds the datatype in the context of a declaration by calling more specific methods
func (getter *DataTypeFactory) declarationBuild() (interface{}, error){
	if len(getter.ctxNode.Children) == 0{
		return getter.declarationSimple()
	}else{
		switch getter.ctxNode.Value.Type {
		case token.ASTERISK:
			return getter.declarationBuildPointer()
		case token.RBRACKET:
			return getter.declarationBuildArray()
		default:
			panic(errorhandler.UnexpectedCompilerError())

		}
	}
}

//declarationSimple return a boolean, a byte or a void depending on the context
func (getter *DataTypeFactory) declarationSimple() (interface{}, error) {
	switch getter.ctxNode.Value.Type {
	case token.TYPEBOOL:
		return symboltable.NewBool(), nil
	case token.TYPEBYTE:
		return symboltable.NewByte(), nil
	case token.VOID:
		return symboltable.NewVoid(), nil
	default:
		panic(errorhandler.UnexpectedCompilerError())

	}
}

//declarationBuildPointer returns a pointer data type. The data type of the reference it points to, is given by moving
//the context to the child of the current node and building its datatype
func (getter *DataTypeFactory) declarationBuildPointer() (interface{}, error) {
	getter.ctxNode = getter.ctxNode.Children[0]
	pointsTo, err := getter.declarationBuild()
	if err != nil {
		return nil, err
	}
	return symboltable.NewPointer(pointsTo), nil
}

//declarationBuildArray validates that the index of an array is valid (by checking its data type) and return a array data type
//the data type of the elements of the array is obtained by moving the context and calling to declarationBuild()
func (getter *DataTypeFactory) declarationBuildArray() (interface{}, error) {
	length := 0
	index := getter.ctxNode.Children[0]
	if index.Value.Type != token.BYTE{
		backup := getter.ctxNode
		getter.ctxNode = index
		dataTypeIndex, err := getter.GetDataType()
		if err != nil{
			return nil, err
		}
		if reflect.TypeOf(dataTypeIndex) != reflect.TypeOf(symboltable.NewByte()){
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.IndexMustBeAByte(line))
			return nil, err
		}
		getter.ctxNode = backup

	}else{
		literal, err := strconv.Atoi(index.Value.Literal)
		if err != nil{
			panic(errorhandler.UnexpectedCompilerError())
		}
		length = literal
	}
	getter.ctxNode = getter.ctxNode.Children[1]
	of, err := getter.declarationBuild()
	if err != nil {
		return nil, err
	}
	return symboltable.NewArray(length, of), nil
}

// GetLeafByRight gets the leaf by walking a tree using the right child of each node.
func GetLeafByRight(head *ast.Node) *ast.Node{
	current := head
	for len(current.Children) != 0{
		current = current.Children[len(current.Children)-1]
	}
	return current
}

// CountNodesToLeafByRight counts how many nodes are left to found
//a leaf by walking a tree using the right child of each node.
func CountNodesToLeafByRight(head *ast.Node) int{
	current := head
	i := 0
	for len(current.Children) != 0{
		i++
		current = current.Children[len(current.Children)-1]
	}
	return i
}


