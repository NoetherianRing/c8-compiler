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
	typebyte := token.NewToken(token.TYPEBYTE, "250", 1)
	typebool := token.NewToken(token.TYPEBOOL, "bool", 1)
	number :=  token.NewToken(token.BYTE, "10", 1)
//	bool :=  token.NewToken(token.BOOL, "true", 1)
//	referenceBool :=  token.NewToken(token.IDENT, "BoolInSymbolTable", 1)
//	referenceByte :=  token.NewToken(token.IDENT, "ByteInSymbolTable", 1)
//	referencePointer :=  token.NewToken(token.IDENT, "PointerInSymbolTable", 1)
//	unknow :=  token.NewToken(token.IDENT, "notStored", 1)

	treePointerToVoid := ast.NewSyntaxTree(ast.NewNode(asterisk))
	treePointerToVoid.Head.AddChild(ast.NewNode(void))

	treeArrayToVoid :=  ast.NewSyntaxTree(ast.NewNode(rbracket))
	treeArrayToVoid.Head.AddChild(ast.NewNode(number))
	treeArrayToVoid.Head.AddChild(ast.NewNode(void))

	testCases := []cases{
		{description: "void",
		ctxNode: ast.NewNode(void),
		expectedDatatype: symboltable.NewVoid(),
		expectedErr: nil,
		},
		{description: "byte",
			ctxNode: ast.NewNode(typebyte),
			expectedDatatype: symboltable.NewByte(),
			expectedErr: nil,
		},
		{description: "bool",
			ctxNode: ast.NewNode(typebool),
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