package semanticAnalyzer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimple(t *testing.T) {
	getter := NewDataTypeGetter()
	scope := symboltable.CreateMainScope()


	type cases struct{
		description     string
		ctxNode         *ast.Node
		expectedDatatype         interface{}
		expectedErr         error

	}

	testCases := []cases{
		{
			description: "byte",
			ctxNode: ast.NewNode(token.NewToken(token.BYTE, "3", 0)),
			expectedDatatype: symboltable.NewByte(),
			expectedErr: nil,
		},
		{
			description: "bool true",
			ctxNode: ast.NewNode(token.NewToken(token.BOOL, "true", 0)),
			expectedDatatype: symboltable.NewBool(),
			expectedErr: nil,
		},

		{
			description: "bool false",
			ctxNode: ast.NewNode(token.NewToken(token.BOOL, "false", 0)),
			expectedDatatype: symboltable.NewBool(),
			expectedErr: nil,
		},

	}

	for _, scenario := range testCases{
		t.Run(scenario.description, func(t *testing.T) {
			getter.SetScope(scope)
			getter.SetCxtNode(scenario.ctxNode)
			datatype, err := getter.simple()
			assert.Equal(t, scenario.expectedErr, err)
			assert.Equal(t, scenario.expectedDatatype, datatype)
		})
	}
}

func TestDeclaration(t *testing.T) {
	getter := NewDataTypeGetter()
	scope := symboltable.CreateMainScope()
	scope.AddSymbol("BoolInSymbolTable", symboltable.NewBool())
	scope.AddSymbol("ByteInSymbolTable", symboltable.NewByte())
	scope.AddSymbol("PointerInSymbolTable", symboltable.NewPointer(nil))
	type cases struct {
		description      string
		ctxNode          *ast.Node
		expectedDatatype interface{}
		expectedErr      error
	}
	asterisk := token.NewToken(token.ASTERISK, token.ASTERISK, 1)
	rbracket := token.NewToken(token.RBRACKET, token.RBRACKET, 1)
	void := token.NewToken(token.VOID, "void", 1)
	typeByte := token.NewToken(token.TYPEBYTE, "byte", 1)
	typeBool := token.NewToken(token.TYPEBOOL, "bool", 1)
	number :=  token.NewToken(token.BYTE, "10", 1)
	negativeNumber :=  token.NewToken(token.BYTE, "-3", 1)
	boolean :=  token.NewToken(token.BOOL, "true", 1)
	referenceByte :=  token.NewToken(token.IDENT, "ByteInSymbolTable", 1)
	unknown :=  token.NewToken(token.IDENT, "notStored", 1)

	datatypeBool := symboltable.NewBool()
	datatypeByte := symboltable.NewByte()
	datatypeVoid := symboltable.NewVoid()


	datatypePointerToBool := 	symboltable.NewPointer(datatypeBool)
	datatypePointerToPointerToBool := 	symboltable.NewPointer(datatypePointerToBool)

	datatypeArrayToByteWithIndex := 	symboltable.NewArray( 10, datatypeByte)
	datatypeArrayToByteWithoutIndex := 	symboltable.NewArray(symboltable.UnknownLength, datatypeByte)
	datatypeArrayToArrayToByte := symboltable.NewArray(symboltable.UnknownLength, datatypeArrayToByteWithIndex)

	datatypeArrayToPointerToBool := symboltable.NewArray(10, datatypePointerToBool)
	datatypePointerToArrayToByte := symboltable.NewPointer(datatypeArrayToByteWithIndex)

	treePointerToVoid := ast.NewSyntaxTree(ast.NewNode(asterisk))
	treePointerToVoid.Head.AddChild(ast.NewNode(void))

	treeArrayToVoid :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToVoid.Head.AddChild(ast.NewNode(number))
	treeArrayToVoid.Head.AddChild(ast.NewNode(void))

	treePointerToBool := ast.NewSyntaxTree(ast.NewNode(asterisk))
	treePointerToBool.Head.AddChild(ast.NewNode(typeBool))

	treeArrayToByteWithIndex :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToByteWithIndex.Head.AddChild(ast.NewNode(number))
	treeArrayToByteWithIndex.Head.AddChild(ast.NewNode(typeByte))

	treeArrayToByteWithoutIndex :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToByteWithoutIndex.Head.AddChild(ast.NewNode(referenceByte))
	treeArrayToByteWithoutIndex.Head.AddChild(ast.NewNode(typeByte))

	treeArrayToByteWithBoolIndex :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToByteWithBoolIndex.Head.AddChild(ast.NewNode(boolean))
	treeArrayToByteWithBoolIndex.Head.AddChild(ast.NewNode(typeByte))

	treeArrayToByteWithUnknownIndex :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToByteWithUnknownIndex.Head.AddChild(ast.NewNode(unknown))
	treeArrayToByteWithUnknownIndex.Head.AddChild(ast.NewNode(typeByte))

	treePointerToPointerToBool := ast.NewSyntaxTree(ast.NewNode(asterisk))
	treePointerToPointerToBool.Head.AddChild(ast.NewNode(asterisk))
	treePointerToPointerToBool.Head.Children[0].AddChild(ast.NewNode(typeBool))

	treeArrayToArrayToByte := ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToArrayToByte.Head.AddChild(ast.NewNode(referenceByte))
	treeArrayToArrayToByte.Head.AddChild(ast.NewNode(rbracket))
	treeArrayToArrayToByte.Head.Children[1].AddChild(ast.NewNode(number))
	treeArrayToArrayToByte.Head.Children[1].AddChild(ast.NewNode(typeByte))

	treeArrayToPointerToBool := ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToPointerToBool.Head.AddChild(ast.NewNode(number))
	treeArrayToPointerToBool.Head.AddChild(ast.NewNode(asterisk))
	treeArrayToPointerToBool.Head.Children[1].AddChild(ast.NewNode(typeBool))

	treePointerToArrayToByte :=  ast.NewSyntaxTree(ast.NewNode(asterisk))
	treePointerToArrayToByte.Head.AddChild(ast.NewNode(rbracket))
	treePointerToArrayToByte.Head.Children[0].AddChild(ast.NewNode(number))
	treePointerToArrayToByte.Head.Children[0].AddChild(ast.NewNode(typeByte))

	treeNegativeIndex := ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeNegativeIndex.Head.AddChild(ast.NewNode(negativeNumber))
	treeNegativeIndex.Head.AddChild(ast.NewNode(typeByte))


	testCases := []cases{
		{description: "void",
		ctxNode: ast.NewNode(void),
		expectedDatatype: symboltable.NewVoid(),
		expectedErr: nil,
		},
		{description: "byte",
			ctxNode: ast.NewNode(typeByte),
			expectedDatatype: symboltable.NewByte(),
			expectedErr: nil,
		},
		{description: "bool",
			ctxNode: ast.NewNode(typeBool),
			expectedDatatype: symboltable.NewBool(),
			expectedErr: nil,
		},
		{
			description: "pointer to void",
			ctxNode: treePointerToVoid.Head,
			expectedDatatype: nil,
			expectedErr: errors.New(errorhandler.PointerToVoid(1)),
		},
		{
			description: "array to void",
			ctxNode: treeArrayToVoid.Head,
			expectedDatatype: nil,
			expectedErr: errors.New(errorhandler.PointerToVoid(1)),
		},
		{
			description: "pointer to bool",
			ctxNode: treePointerToBool.Head,
			expectedDatatype: datatypePointerToBool,
			expectedErr: nil,
		},
		{
			description:      "array to number with index",
			ctxNode:          treeArrayToByteWithIndex.Head,
			expectedDatatype: datatypeArrayToByteWithIndex,
			expectedErr:      nil,
		},
		{
			description:      "array to number without index",
			ctxNode:          treeArrayToByteWithoutIndex.Head,
			expectedDatatype: datatypeArrayToByteWithoutIndex,
			expectedErr:      nil,
		},
		{
			description:      "array to number with boolean index",
			ctxNode:          treeArrayToByteWithBoolIndex.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.IndexMustBeAByte(1)),
		},
		{
			description:      "array to number with unknown index",
			ctxNode:          treeArrayToByteWithUnknownIndex.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.UnresolvedReference(1, unknown.Literal)),
		},
		{
			description: "pointer to pointer to bool",
			ctxNode: treePointerToPointerToBool.Head,
			expectedDatatype: datatypePointerToPointerToBool,
			expectedErr: nil,
		},
		{
			description: "array to array to byte",
			ctxNode: treeArrayToArrayToByte.Head,
			expectedDatatype: datatypeArrayToArrayToByte,
			expectedErr: nil,
		},
		{
			description: "array to pointer to bool",
			ctxNode: treeArrayToPointerToBool.Head,
			expectedDatatype: datatypeArrayToPointerToBool,
			expectedErr: nil,
		},
		{
			description: "pointer to array to byte",
			ctxNode: treePointerToArrayToByte.Head,
			expectedDatatype: datatypePointerToArrayToByte,
			expectedErr: nil,
		},
		{
			description: "byte",
			ctxNode: ast.NewNode(typeByte),
			expectedDatatype: datatypeByte,
			expectedErr: nil,
		},
		{
			description: "bool",
			ctxNode: ast.NewNode(typeBool),
			expectedDatatype: datatypeBool,
			expectedErr: nil,
		},
		{
			description: "void",
			ctxNode: ast.NewNode(void),
			expectedDatatype: datatypeVoid,
			expectedErr: nil,
		},
		{
			description:      "negative index",
			ctxNode:          treeNegativeIndex.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NegativeIndex(1)),
		},
	}
	for _, scenario := range testCases{
		t.Run(scenario.description, func(t *testing.T) {
			getter.SetScope(scope)
			getter.SetCxtNode(scenario.ctxNode)
			datatype, err := getter.declaration()
			assert.Equal(t, scenario.expectedErr, err)
			assert.Equal(t, scenario.expectedDatatype, datatype)
		})
	}
}

func TestReference(t *testing.T){
	getter := NewDataTypeGetter()
	scope := symboltable.CreateMainScope()

	type cases struct {
		description      string
		ctxNode          *ast.Node
		expectedDatatype interface{}
		expectedErr      error
	}

	datatypeBool := symboltable.NewBool()
	datatypeByte := symboltable.NewByte()
	datatypeVoid := symboltable.NewVoid()

	voidFunc := symboltable.NewFunction(datatypeVoid, nil)
	datatypePointerToBool := 	symboltable.NewPointer(datatypeBool)
	datatypePointerToByte := 	symboltable.NewPointer(datatypeByte)
	datatypePointerToPointerToBool := 	symboltable.NewPointer(datatypePointerToBool)

	datatypeArrayToByteWithIndex := 	symboltable.NewArray( 10, datatypeByte)
	datatypeArrayToByteWithoutIndex := 	symboltable.NewArray(symboltable.UnknownLength, datatypeByte)
	datatypeArrayToArrayToByte := symboltable.NewArray(symboltable.UnknownLength, datatypeArrayToByteWithIndex)

	datatypeArrayToPointerToBool := symboltable.NewArray(10, datatypePointerToBool)
	datatypePointerToArrayToByte := symboltable.NewPointer(datatypeArrayToByteWithIndex)

	varBool := token.NewToken(token.IDENT, "varBool", 1)
	varByte := token.NewToken(token.IDENT, "varByte", 1)
	myFunc := token.NewToken(token.IDENT, "myFunc", 1)
	pointerToBool := token.NewToken(token.IDENT, "pointerToBool", 1)
	pointerToByte := token.NewToken(token.IDENT, "pointerToByte", 1)
	pointerToPointerToBool := token.NewToken(token.IDENT, "pointerToPointerToBool", 1)
	arrayToByteWithIndex := token.NewToken(token.IDENT, "arrayToByteWithIndex", 1)
	arrayToByteWithoutIndex := token.NewToken(token.IDENT, "arrayToByteWithoutIndex", 1)
	arrayToArrayToByte := token.NewToken(token.IDENT, "arrayToArrayToByte", 1)
	arrayToPointerToBool := token.NewToken(token.IDENT, "arrayToPointerToBool", 1)
	pointerToArrayToByte := token.NewToken(token.IDENT, "pointerToArrayToByte", 1)
	unknown  := token.NewToken(token.IDENT, "unknown", 1)

	scope.AddSymbol(varBool.Literal, datatypeBool)
	scope.AddSymbol(varByte.Literal, datatypeByte)
	scope.AddSymbol(myFunc.Literal, voidFunc)
	scope.AddSymbol(pointerToBool.Literal, datatypePointerToBool)
	scope.AddSymbol(pointerToByte.Literal, datatypePointerToByte)
	scope.AddSymbol(pointerToPointerToBool.Literal, datatypePointerToPointerToBool)
	scope.AddSymbol(arrayToByteWithIndex.Literal, datatypeArrayToByteWithIndex)
	scope.AddSymbol(arrayToByteWithoutIndex.Literal, datatypeArrayToByteWithoutIndex)
	scope.AddSymbol(arrayToArrayToByte.Literal, datatypeArrayToArrayToByte)
	scope.AddSymbol(arrayToPointerToBool.Literal, datatypeArrayToPointerToBool)
	scope.AddSymbol(pointerToArrayToByte.Literal, datatypePointerToArrayToByte)

	testCases := []cases{
		{
			description:      "Unresolved reference",
			ctxNode:          ast.NewNode(unknown),
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.UnresolvedReference(unknown.Line, unknown.Literal)),
		},
		{
			description:      "Identifier is a function",
			ctxNode:          ast.NewNode(myFunc),
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.IdentifierIsFunction(myFunc.Line, myFunc.Literal)),
		},
		{
			description:      "varBool",
			ctxNode:          ast.NewNode(varBool),
			expectedDatatype: datatypeBool,
			expectedErr:      nil,
		},
		{
			description:      "varByte",
			ctxNode:          ast.NewNode(varByte),
			expectedDatatype: datatypeByte,
			expectedErr:      nil,
		},
		{
			description:      "pointerToBool",
			ctxNode:          ast.NewNode(pointerToBool),
			expectedDatatype: datatypePointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "pointerToByte",
			ctxNode:          ast.NewNode(pointerToByte),
			expectedDatatype: datatypePointerToByte,
			expectedErr:      nil,
		},
		{
			description:      "pointerToPointerToBool",
			ctxNode:          ast.NewNode(pointerToPointerToBool),
			expectedDatatype: datatypePointerToPointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "arrayToByteWithIndex",
			ctxNode:          ast.NewNode(arrayToByteWithIndex),
			expectedDatatype: datatypeArrayToByteWithIndex,
			expectedErr:      nil,
		},
		{
			description:      "arrayToByteWithoutIndex",
			ctxNode:          ast.NewNode(arrayToByteWithoutIndex),
			expectedDatatype: datatypeArrayToByteWithoutIndex,
			expectedErr:      nil,
		},
		{
			description:      "arrayToArrayToByte",
			ctxNode:          ast.NewNode(arrayToArrayToByte),
			expectedDatatype: datatypeArrayToArrayToByte,
			expectedErr:      nil,
		},
		{
			description:      "arrayToPointerToBool",
			ctxNode:          ast.NewNode(arrayToPointerToBool),
			expectedDatatype: datatypeArrayToPointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "pointerToArrayToByte",
			ctxNode:          ast.NewNode(pointerToArrayToByte),
			expectedDatatype: datatypePointerToArrayToByte,
			expectedErr:      nil,
		},
	}
	for _, scenario := range testCases{
		t.Run(scenario.description, func(t *testing.T) {
			getter.SetScope(scope)
			getter.SetCxtNode(scenario.ctxNode)
			datatype, err := getter.reference()
			assert.Equal(t, scenario.expectedErr, err)
			assert.Equal(t, scenario.expectedDatatype, datatype)
		})
	}
}

func TestFunctionCall(t *testing.T){
	getter := NewDataTypeGetter()
	scope := symboltable.CreateMainScope()

	type cases struct {
		description      string
		ctxNode          *ast.Node
		expectedDatatype interface{}
		expectedErr      error
	}

	datatypeBool := symboltable.NewBool()
	datatypeByte := symboltable.NewByte()
	datatypeVoid := symboltable.NewVoid()

	datatypePointerToBool := 	symboltable.NewPointer(datatypeBool)
	datatypePointerToPointerToBool := 	symboltable.NewPointer(datatypePointerToBool)

	datatypeArrayToByteWithIndex := 	symboltable.NewArray( 10, datatypeByte)
	datatypeArrayToArrayToByte := symboltable.NewArray(0, datatypeArrayToByteWithIndex)

	datatypePointerToArrayToByte := symboltable.NewPointer(datatypeArrayToByteWithIndex)


	datatypeFuncVoidWithoutParameters := symboltable.NewFunction(datatypeVoid, nil)
	funcVoidWithoutParameters := token.NewToken(token.IDENT, "funcVoidWithoutParameters", 1)
	scope.AddSymbol(funcVoidWithoutParameters.Literal, datatypeFuncVoidWithoutParameters)

	datatypeFuncByteWithoutParameters := symboltable.NewFunction(datatypeByte, nil)
	funcByteWithoutParameters := token.NewToken(token.IDENT, "funcByteWithoutParameters", 1)
	scope.AddSymbol(funcByteWithoutParameters.Literal, datatypeFuncByteWithoutParameters)

	datatypeFuncBoolWithoutParameters := symboltable.NewFunction(datatypeBool, nil)
	funcBoolWithoutParameters := token.NewToken(token.IDENT, "funcBoolWithoutParameters", 1)
	scope.AddSymbol(funcBoolWithoutParameters.Literal, datatypeFuncBoolWithoutParameters)


	args := make([]interface{},0)
	args = append(args,  datatypeByte)
	datatypeFuncVoidParamByte := symboltable.NewFunction(datatypeVoid, args)
	funcVoidParamByte := token.NewToken(token.IDENT, "funcVoidParamByte", 1)
	scope.AddSymbol(funcVoidParamByte.Literal, datatypeFuncVoidParamByte)


	args = make([]interface{},0)
	args = append(args,  datatypePointerToPointerToBool)
	datatypeFuncByteParamPointerToPointerToBool := symboltable.NewFunction(datatypeByte, args)
	funcByteParamPointerToPointerToBool := token.NewToken(token.IDENT, "funcByteParamPointerToPointerToBool", 1)
	scope.AddSymbol(funcByteParamPointerToPointerToBool.Literal, datatypeFuncByteParamPointerToPointerToBool)

	args = make([]interface{},0)
	args = append(args,  datatypeByte)
	args = append(args,  datatypePointerToPointerToBool)
	datatypeFuncBoolParamByteAndPointerToPointerToBool := symboltable.NewFunction(datatypeBool, args)
	funcBoolParamByteAndPointerToPointerToBool := token.NewToken(token.IDENT, "funcBoolParamByteAndPointerToPointerToBool", 1)
	scope.AddSymbol(funcBoolParamByteAndPointerToPointerToBool.Literal, datatypeFuncBoolParamByteAndPointerToPointerToBool)

	args = make([]interface{},0)
	args = append(args,  datatypePointerToArrayToByte)

	datatypeFuncArrayToArrayToByteParamPointerToArrayToByte := symboltable.NewFunction(datatypeArrayToArrayToByte, args)
	funcArrayToArrayToByteParamPointerToArrayToByte := token.NewToken(token.IDENT, "funcArrayToArrayToByteParamPointerToArrayToByte", 1)
	scope.AddSymbol(funcArrayToArrayToByteParamPointerToArrayToByte.Literal, datatypeFuncArrayToArrayToByteParamPointerToArrayToByte)

	parentheses := token.NewToken(token.RPAREN, token.RPAREN, 1)
	comma := token.NewToken(token.COMMA, token.COMMA, 1)
	number := token.NewToken(token.BYTE, "8", 1)
	unknown := token.NewToken(token.IDENT, "unknown", 1)
	identifierPointerToPointerToBool := token.NewToken(token.IDENT, "var1", 1)
	scope.AddSymbol(identifierPointerToPointerToBool.Literal, datatypePointerToPointerToBool)

	identifierPointerToArrayToByte := token.NewToken(token.IDENT, "var2", 1)
	scope.AddSymbol(identifierPointerToArrayToByte.Literal, datatypePointerToArrayToByte)

	treeVoidWithoutParams := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeVoidWithoutParams.Head.AddChild(ast.NewNode(funcVoidWithoutParameters))

	treeByteWithoutParams := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeByteWithoutParams.Head.AddChild(ast.NewNode(funcByteWithoutParameters))

	treeBoolWithoutParams := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeBoolWithoutParams.Head.AddChild(ast.NewNode(funcBoolWithoutParameters))

	treeVoidParamByte := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeVoidParamByte.Head.AddChild(ast.NewNode(funcVoidParamByte))
	treeVoidParamByte.Head.AddChild(ast.NewNode(number))

	treeByteParamPointerToPointerToBool := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeByteParamPointerToPointerToBool.Head.AddChild(ast.NewNode(funcByteParamPointerToPointerToBool))
	treeByteParamPointerToPointerToBool.Head.AddChild(ast.NewNode(identifierPointerToPointerToBool))

	treeBoolParamByteAndPointerToPointerToBool := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeBoolParamByteAndPointerToPointerToBool.Head.AddChild(ast.NewNode(funcBoolParamByteAndPointerToPointerToBool))
	treeBoolParamByteAndPointerToPointerToBool.Head.AddChild(ast.NewNode(comma))
	treeBoolParamByteAndPointerToPointerToBool.Head.Children[1].AddChild(ast.NewNode(number))
	treeBoolParamByteAndPointerToPointerToBool.Head.Children[1].AddChild(ast.NewNode(identifierPointerToPointerToBool))


	treeArrayToArrayToByteParamPointerToArrayToByte := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeArrayToArrayToByteParamPointerToArrayToByte.Head.AddChild(ast.NewNode(funcArrayToArrayToByteParamPointerToArrayToByte))
	treeArrayToArrayToByteParamPointerToArrayToByte.Head.AddChild(ast.NewNode(identifierPointerToArrayToByte))

	treeUnresolvedReference :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeUnresolvedReference.Head.AddChild(ast.NewNode(unknown))

	treeIsNotAFunction := ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeIsNotAFunction.Head.AddChild(ast.NewNode(identifierPointerToArrayToByte))

	treeMoreNumberOfParams :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeMoreNumberOfParams.Head.AddChild(ast.NewNode(funcBoolWithoutParameters))
	treeMoreNumberOfParams.Head.AddChild(ast.NewNode(comma))
	treeMoreNumberOfParams.Head.Children[1].AddChild(ast.NewNode(number))
	treeMoreNumberOfParams.Head.Children[1].AddChild(ast.NewNode(number))

	treeMoreNumberOfParams2 :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeMoreNumberOfParams2.Head.AddChild(ast.NewNode(funcBoolParamByteAndPointerToPointerToBool))
	treeMoreNumberOfParams2.Head.AddChild(ast.NewNode(comma))
	treeMoreNumberOfParams2.Head.Children[1].AddChild(ast.NewNode(number))
	treeMoreNumberOfParams2.Head.Children[1].AddChild(ast.NewNode(comma))
	treeMoreNumberOfParams2.Head.Children[1].Children[1].AddChild(ast.NewNode(identifierPointerToPointerToBool))
	treeMoreNumberOfParams2.Head.Children[1].Children[1].AddChild(ast.NewNode(number))

	treeLessNumberOfParams :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeLessNumberOfParams.Head.AddChild(ast.NewNode(funcVoidParamByte))

	treeLessNumberOfParams2 :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeLessNumberOfParams2.Head.AddChild(ast.NewNode(funcBoolParamByteAndPointerToPointerToBool))
	treeLessNumberOfParams2.Head.AddChild(ast.NewNode(number))


	treeDatatypeMismatches1 :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeDatatypeMismatches1.Head.AddChild(ast.NewNode(funcVoidParamByte))
	treeDatatypeMismatches1.Head.AddChild(ast.NewNode(identifierPointerToArrayToByte))

	treeDatatypeMismatches2 :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeDatatypeMismatches2.Head.AddChild(ast.NewNode(funcBoolParamByteAndPointerToPointerToBool))
	treeDatatypeMismatches2.Head.AddChild(ast.NewNode(comma))
	treeDatatypeMismatches2.Head.Children[1].AddChild(ast.NewNode(identifierPointerToPointerToBool))
	treeDatatypeMismatches2.Head.Children[1].AddChild(ast.NewNode(identifierPointerToPointerToBool))

	treeDatatypeMismatches3 :=  ast.NewSyntaxTree(ast.NewNode(parentheses))
	treeDatatypeMismatches3.Head.AddChild(ast.NewNode(funcBoolParamByteAndPointerToPointerToBool))
	treeDatatypeMismatches3.Head.AddChild(ast.NewNode(comma))
	treeDatatypeMismatches3.Head.Children[1].AddChild(ast.NewNode(number))
	treeDatatypeMismatches3.Head.Children[1].AddChild(ast.NewNode(number))




	testCases := []cases{
		{
			description:      "void without params",
			ctxNode:          treeVoidWithoutParams.Head,
			expectedDatatype: datatypeVoid,
			expectedErr:      nil,
		},
		{
			description:      "byte without params",
			ctxNode:          treeByteWithoutParams.Head,
			expectedDatatype: datatypeByte,
			expectedErr:      nil,
		},
		{
			description:      "bool without params",
			ctxNode:          treeBoolWithoutParams.Head,
			expectedDatatype: datatypeBool,
			expectedErr:      nil,
		},
		{
			description:      "void param byte",
			ctxNode:          treeVoidParamByte.Head,
			expectedDatatype: datatypeVoid,
			expectedErr:      nil,
		},
		{
			description:      "byte param pointer to pointer to bool",
			ctxNode:          treeByteParamPointerToPointerToBool.Head,
			expectedDatatype: datatypeByte,
			expectedErr:      nil,
		},
		{
			description:      "bool param byte, pointer to pointer to bool",
			ctxNode:          treeBoolParamByteAndPointerToPointerToBool.Head,
			expectedDatatype: datatypeBool,
			expectedErr:      nil,
		},
		{
			description:      "Array to Array to Byte Param Pointer to Array to Byte",
			ctxNode:          treeArrayToArrayToByteParamPointerToArrayToByte.Head,
			expectedDatatype: datatypeArrayToArrayToByte,
			expectedErr:      nil,
		},
		{
			description:      "Unresolved Reference",
			ctxNode:          treeUnresolvedReference.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.UnresolvedReference(unknown.Line, unknown.Literal)),
		},
		{
			description:      "Identifier is Not A Function",
			ctxNode:          treeIsNotAFunction.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.IdentifierIsNotFunction(identifierPointerToArrayToByte.Line, identifierPointerToArrayToByte.Literal)),
		},
		{
			description:      "Number Of Parameters Doesnt Match 2 = 0",
			ctxNode:          treeMoreNumberOfParams.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NumberOfParametersDoesntMatch(1, 2, 0)),
		},
		{
			description:      "Number Of Parameters Doesnt Match 3 = 2",
			ctxNode:          treeMoreNumberOfParams2.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NumberOfParametersDoesntMatch(1, 3, 2)),
		},
		{
			description:      "Number Of Parameters Doesnt Match 0 = 1",
			ctxNode:          treeLessNumberOfParams.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NumberOfParametersDoesntMatch(1, 0, 1)),
		},
		{
			description:      "Number Of Parameters Doesnt Match 1 = 2",
			ctxNode:          treeLessNumberOfParams2.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NumberOfParametersDoesntMatch(1, 1, 2)),
		},
		{
			description:      "Data type mismatches *[]bool = byte",
			ctxNode:          treeDatatypeMismatches1.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.DataTypesMismatch(1,symboltable.Fmt(datatypePointerToArrayToByte), "=", symboltable.Fmt(datatypeByte))),
		},
		{
			description:      "Data type mismatches **bool = byte",
			ctxNode:          treeDatatypeMismatches2.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.DataTypesMismatch(1,symboltable.Fmt(datatypePointerToPointerToBool), "=", symboltable.Fmt(datatypeByte))),
		},
		{
			description:      "Data type mismatches byte = **bool",
			ctxNode:          treeDatatypeMismatches3.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.DataTypesMismatch(1, symboltable.Fmt(datatypeByte), "=",symboltable.Fmt(datatypePointerToPointerToBool))),
		},

	}
	for _, scenario := range testCases{
		t.Run(scenario.description, func(t *testing.T) {
			getter.SetScope(scope)
			getter.SetCxtNode(scenario.ctxNode)
			datatype, err := getter.functionCall()
			assert.Equal(t, scenario.expectedErr, err)
			assert.Equal(t, scenario.expectedDatatype, datatype)
		})
	}
}

func TestDereference(t *testing.T){

	getter := NewDataTypeGetter()
	scope := symboltable.CreateMainScope()


	type cases struct{
		description     string
		ctxNode         *ast.Node
		expectedDatatype         interface{}
		expectedErr         error

	}
	asterisk := token.NewToken(token.ASTERISK, token.ASTERISK, 1)
	rbracket := token.NewToken(token.RBRACKET, token.RBRACKET, 1)
//	void := token.NewToken(token.VOID, "void", 1)
	number :=  token.NewToken(token.BYTE, "8", 1)
	negativeNumber :=  token.NewToken(token.BYTE, "-100", 1)
	numberOutOfBound :=  token.NewToken(token.BYTE, "100", 1)
//	boolean :=  token.NewToken(token.BOOL, "true", 1)
//	referenceByte :=  token.NewToken(token.IDENT, "ByteInSymbolTable", 1)
	unknown :=  token.NewToken(token.IDENT, "notStored", 1)

	datatypeBool := symboltable.NewBool()
	datatypeByte := symboltable.NewByte()


	datatypePointerToBool := 	symboltable.NewPointer(datatypeBool)
	datatypePointerToPointerToBool := 	symboltable.NewPointer(datatypePointerToBool)

	identPPBool := token.NewToken(token.IDENT, "PointerPointerBool", 1)
	scope.AddSymbol(identPPBool.Literal, datatypePointerToPointerToBool)

	tree1 := ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree1.Head.AddChild(ast.NewNode(asterisk))
	tree1.Head.Children[0].AddChild(ast.NewNode(identPPBool))

	tree2 := ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree2.Head.AddChild(ast.NewNode(identPPBool))

	tree3 := ast.NewSyntaxTree(ast.NewNode(identPPBool))

	datatypeArrayToByteWithIndex := 	symboltable.NewArray( 10, datatypeByte)
	datatypeArrayToArrayToByte := symboltable.NewArray(symboltable.UnknownLength, datatypeArrayToByteWithIndex)

	identAAByte := token.NewToken(token.IDENT, "ArrayArrayByte", 1)
	scope.AddSymbol(identAAByte.Literal, datatypeArrayToArrayToByte)

	tree4:= ast.NewSyntaxTree(ast.NewNode(rbracket))
	tree4.Head.AddChild(ast.NewNode(number))
	tree4.Head.AddChild(ast.NewNode(identAAByte))

	datatypeArrayToPointerToBool := symboltable.NewArray(10, datatypePointerToBool)

	identAPBool := token.NewToken(token.IDENT, "ArrayPointerBool", 1)
	scope.AddSymbol(identAPBool.Literal, datatypeArrayToPointerToBool)

	tree5:= ast.NewSyntaxTree(ast.NewNode(rbracket))
	tree5.Head.AddChild(ast.NewNode(number))
	tree5.Head.AddChild(ast.NewNode(identAPBool))

	datatypePointerToArrayToByte := symboltable.NewPointer(datatypeArrayToByteWithIndex)
	identPAByte := token.NewToken(token.IDENT, "PointerArrayByte", 1)
	scope.AddSymbol(identPAByte.Literal, datatypePointerToArrayToByte)

	tree6:= ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree6.Head.AddChild(ast.NewNode(identPAByte))

	tree7:= ast.NewSyntaxTree(ast.NewNode(rbracket))
	tree7.Head.AddChild(ast.NewNode(numberOutOfBound))
	tree7.Head.AddChild(ast.NewNode(identAPBool))

	tree8:= ast.NewSyntaxTree(ast.NewNode(rbracket))
	tree8.Head.AddChild(ast.NewNode(negativeNumber))
	tree8.Head.AddChild(ast.NewNode(identAPBool))

	tree9:= ast.NewSyntaxTree(ast.NewNode(rbracket))
	tree9.Head.AddChild(ast.NewNode(identAPBool))
	tree9.Head.AddChild(ast.NewNode(identAPBool))

	tree10:= ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree10.Head.AddChild(ast.NewNode(unknown))

	tree11:= ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree11.Head.AddChild(ast.NewNode(identAPBool))

	tree12:= ast.NewSyntaxTree(ast.NewNode(asterisk))
	tree12.Head.AddChild(ast.NewNode(asterisk))
	tree12.Head.Children[0].AddChild(ast.NewNode(asterisk))
	tree12.Head.Children[0].Children[0].AddChild(ast.NewNode(identPPBool))

	testCases := []cases{
		{
			description:      "**identPPBool",
			ctxNode:          tree1.Head,
			expectedDatatype: datatypeBool,
			expectedErr:      nil,
		},
		{
			description:      "*identPPBool",
			ctxNode:          tree2.Head,
			expectedDatatype: datatypePointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "identPPBool",
			ctxNode:          tree3.Head,
			expectedDatatype: datatypePointerToPointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "[10]identAAByte",
			ctxNode:          tree4.Head,
			expectedDatatype: datatypeArrayToByteWithIndex,
			expectedErr:      nil,
		},
		{
			description:      "[8]identAPBool",
			ctxNode:          tree5.Head,
			expectedDatatype: datatypePointerToBool,
			expectedErr:      nil,
		},
		{
			description:      "*identPAByte",
			ctxNode:          tree6.Head,
			expectedDatatype: datatypeArrayToByteWithIndex,
			expectedErr:      nil,
		},
		{
			description:      "[100]identAPBool",
			ctxNode:          tree7.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.IndexOutOfBounds(identAPBool.Line)),
		},
		{
			description:      "[-100]identAPBool",
			ctxNode:          tree8.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.NegativeIndex(identAPBool.Line)),
		},
		{
			description:      "[identAPBool]identAPBool",
			ctxNode:          tree9.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.IndexMustBeAByte(identAPBool.Line)),
		},
		{
			description:      "Unresolved Reference",
			ctxNode:          tree10.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.UnresolvedReference(unknown.Line, unknown.Literal)),
		},
		{
			description:      "InvalidIndirectOf identAPBool",
			ctxNode:          tree11.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.InvalidIndirectOf(identAPBool.Line, identAPBool.Literal)),
		},
		{
			description:      "InvalidIndirectOf identPPBool",
			ctxNode:          tree12.Head,
			expectedDatatype: nil,
			expectedErr:      errors.New(errorhandler.InvalidIndirectOf(identPPBool.Line, identPPBool.Literal)),
		},
	}
	for _, scenario := range testCases{
		t.Run(scenario.description, func(t *testing.T) {
			getter.SetScope(scope)
			getter.SetCxtNode(scenario.ctxNode)
			datatype, err := getter.dereference()
			assert.Equal(t, scenario.expectedErr, err)
			assert.Equal(t, scenario.expectedDatatype, datatype)
		})
	}
}