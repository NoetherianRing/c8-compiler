package parser

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/token"
	"strconv"
)

type SemanticAnalyzer struct{
	symbolTable   map[string]token.Symbol
	checkSemantic map[token.Type]func(tree *ast.SyntaxTree, scope string) error
	getDataType   map[token.Type]func(tree *ast.SyntaxTree, scope string) (*token.DataType, error)
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	table := make(map[string]token.Symbol)
	checkSemantic := make(map[token.Type]func(tree *ast.SyntaxTree, scope string) error)
	getDataType := make(map[token.Type]func(tree *ast.SyntaxTree, scope string) (*token.DataType, error))

	analyzer := new(SemanticAnalyzer)

	checkSemantic[token.RBRACE] = analyzer.block
	checkSemantic[token.LET] = analyzer.let

	getDataType[token.LAND] = analyzer.getLogicOperator
	getDataType[token.LOR] = analyzer.getLogicOperator
	getDataType[token.EQEQ] = analyzer.getEQEQ
	getDataType[token.NOTEQ] = analyzer.getEQEQ
	getDataType[token.LT] = analyzer.getComparison
	getDataType[token.LTEQ] = analyzer.getComparison
	getDataType[token.GT] = analyzer.getComparison
	getDataType[token.GTEQ] = analyzer.getComparison
	getDataType[token.BYTE] = analyzer.getByte
	getDataType[token.BOOL] = analyzer.getBool
	getDataType[token.BANG] = analyzer.getNOT
	getDataType[token.DOLLAR] = analyzer.getAddress
	getDataType[token.TYPEBOOL] = analyzer.getBool
	getDataType[token.TYPEBYTE] = analyzer.getByte
	getDataType[token.VOID] = analyzer.getVoid

	getDataType[token.ASTERISK] = analyzer.getAsterisk
	getDataType[token.IDENT] = analyzer.getIdent
	getDataType[token.RBRACKET] = analyzer.getArray
	analyzer.getDataType = getDataType

	analyzer.checkSemantic = checkSemantic

	analyzer.symbolTable = table
	return analyzer
}


func (analyzer *SemanticAnalyzer) block(tree *ast.SyntaxTree, scope string) error {
	subscopes := 0
	for _, child := range tree.Current.Children{
		newScope := scope
		tree.Current = child
		if child.Value.Type == token.IF || child.Value.Type == token.ELSE ||
			child.Value.Type == token.WHILE || child.Value.Type == token.FUNCTION{
			subscopes++
			newScope = newScope + "_" + strconv.Itoa(subscopes)
 		}
	    err := analyzer.checkSemantic[child.Value.Type](tree, newScope)
		if err != nil{
			return err
		}
	}
	return nil
}

func (analyzer *SemanticAnalyzer) let(tree *ast.SyntaxTree, scope string) error {
	newVar := scope + "_" +tree.Current.Children[0].Value.Literal
	_, exist := analyzer.symbolTable[newVar]
	if exist{
		errorString := "semantic errorhandler\nin line: "+ strconv.Itoa(tree.Current.Children[0].Value.Line) +
			"\nthe variable name " +tree.Current.Children[0].Value.Literal +" its already taken in the scope\n"

		return errors.New(errorString)
	}else{


	}

	return nil
}

func (analyzer *SemanticAnalyzer) GetSymbolTable() token.SymbolTable{
	return analyzer.symbolTable
}




/*
func getType(node *ast.Node) (string, errorhandler){
	switch node.value.Literal {
	case token.LAND:
		return token.TYPEBOOL, nil

	}
}
*/