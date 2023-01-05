package parser

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/token"
	"strconv"
)

//In this file there are a set of functions that are used for the semantic analyser to get the data type of certain operators (if they are valid).

//getLogicOperator is used to inform the data type of expressions with "||" and "&&" as operators of higher precedence
//as well as to validate that sub-expressions on the left and right side of the operator are of a valid data type.
func (analyzer *SemanticAnalyzer) getLogicOperator(tree *ast.SyntaxTree, scope string) (*token.DataType, error){
	boolDataType := token.NewDataType(token.DataTypeBool, 1, 0, nil)
	current := tree.Current

	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(current.Children) != 2{
		panic(errorhandler.UnexpectedCompilerError())
	}
	tree.Current = current.Children[0]
	leftDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil {
		return nil, err
	}
	if leftDataType.Kind != token.DataTypeBool{
		errorString := errorhandler.UnexpectedDataType(current.Value.Line, boolDataType.Fmt(), leftDataType.Fmt())
		return nil, errors.New(errorString)
	}

	tree.Current = current.Children[1]
	rightDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)
	if err != nil{
		return nil, err
	}

	if rightDataType.Kind != token.DataTypeBool{
		errorString := errorhandler.UnexpectedDataType(current.Value.Line, boolDataType.Fmt(), rightDataType.Fmt())
		return nil, errors.New(errorString)
	}

	return boolDataType, nil

}

//getEQEQ is used to inform the data type of expressions with "==" and "!=" as operators of higher precedence
//as well as to validate that sub-expressions on the left and right side of the operator are of a valid data type.
func (analyzer *SemanticAnalyzer) getEQEQ(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	boolDataType := token.NewDataType(token.DataTypeBool, 1, 0, nil)
	current := tree.Current
	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(current.Children) != 2{
		panic(errorhandler.UnexpectedCompilerError())
	}
	tree.Current = current.Children[0]
	leftDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil {
		return nil, err
	}
	tree.Current = current.Children[1]
	rightDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)
	if err != nil{
		return nil, err
	}
	if leftDataType.Kind != rightDataType.Kind{
		return nil, errors.New(errorhandler.DataTypesDontMatch(current.Value.Line, leftDataType.Fmt(),
			current.Value.Literal, rightDataType.Fmt()))
	}
	return boolDataType, nil
}

//getComparison is used to inform the data type of expressions with "<", "<=", ">" and ">=" as operators of higher precedence
//as well as to validate that sub-expressions on the left and right side of the operator are of a valid data type.
func (analyzer *SemanticAnalyzer) getComparison(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	boolDataType := token.NewDataType(token.DataTypeBool, 1, 0, nil)
	current := tree.Current
	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(current.Children) != 2{
		panic(errorhandler.UnexpectedCompilerError())
	}
	tree.Current = current.Children[0]
	leftDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil {
		return nil, err
	}
	if leftDataType.Kind != token.DataTypeByte && leftDataType.Kind != token.DataTypePointer{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, "numeric", leftDataType.Fmt()))
	}
	tree.Current = current.Children[1]
	rightDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)
	if err != nil{
		return nil, err
	}
	if rightDataType.Kind != token.DataTypeByte && rightDataType.Kind != token.DataTypePointer{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, "numeric", rightDataType.Fmt()))
	}
	if leftDataType != rightDataType{
		return nil, errors.New(errorhandler.DataTypesDontMatch(current.Value.Line, leftDataType.Fmt(),
			current.Value.Literal, rightDataType.Fmt()))
	}
	return boolDataType, nil

}

//getOperatorLeftGreaterThanRightSize is used to inform the data type of expressions with
//">>", "<<", "+", "-", "/", "%" and "*" (when * is being used as a multiplication symbol)
//as operators of higher precedence as well as to validate that sub-expressions on the left and right side of the operator are of a valid data type.
//Specifically, the size of the operand on the right must be higher than the operand on the left.
func (analyzer *SemanticAnalyzer) getOperatorLeftGreaterThanRightSize(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	current := tree.Current
	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(current.Children) != 2{
		panic(errorhandler.UnexpectedCompilerError())
	}

	tree.Current = current.Children[0]
	leftDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil {
		return nil, err
	}

	if leftDataType.Kind != token.DataTypeByte && leftDataType.Kind != token.DataTypePointer{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, "numeric", leftDataType.Fmt()))
	}

	tree.Current = current.Children[1]
	rightDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil{
		return nil, err
	}

	if rightDataType.Kind != token.DataTypeByte{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line,
			token.NewDataType(token.DataTypeByte, 1, 0, nil).Fmt(), rightDataType.Fmt()))
	}

	return leftDataType, nil
}

//getOperatorEqualChildrenSize is used to inform the data type of expressions with
//"&", "|", "^" as operators of higher precedence as well as to validate that sub-expressions on the left and right side of the operator are of a valid data type.
//Specifically, the size of the operand on the right must be equal to the operand on the left.
func (analyzer *SemanticAnalyzer) getOperatorEqualChildrenSize(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	current := tree.Current

	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(current.Children) != 2{
		panic(errorhandler.UnexpectedCompilerError())
	}

	tree.Current = current.Children[0]
	leftDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)


	if err != nil {
		return nil, err
	}

	if leftDataType.Kind != token.DataTypeByte && leftDataType.Kind != token.DataTypePointer{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, "numeric", leftDataType.Fmt()))
	}

	tree.Current = current.Children[1]
	rightDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil{
		return nil, err
	}

	if rightDataType.Kind != token.DataTypeByte && rightDataType.Kind != token.DataTypePointer{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, "numeric", rightDataType.Fmt()))
	}

	if leftDataType != rightDataType{
		return nil, errors.New(errorhandler.DataTypesDontMatch(current.Value.Line, leftDataType.Fmt(),
			current.Value.Literal, rightDataType.Fmt()))
	}
	return leftDataType, nil
}

//getNOT is used to inform the data type of expressions with
//"!" as operators of higher precedence as well as to validate that the sub-expression on the right is of a valid data type.
func (analyzer *SemanticAnalyzer) getNOT(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	boolDataType := token.NewDataType(token.DataTypeBool, 1, 0, nil)

	current := tree.Current

	//Due to the syntax analyzer, the tree should not allow the current node to have anything but one child.
	if len(current.Children) != 1{
		panic(errorhandler.UnexpectedCompilerError())
	}

	tree.Current = current.Children[0]

	childDataType, err := analyzer.getDataType[tree.Current.Value.Type](tree, scope)

	if err != nil {
		return nil, err
	}

	if childDataType.Kind != token.DataTypeByte{
		return nil, errors.New(errorhandler.UnexpectedDataType(tree.Current.Value.Line, boolDataType.Fmt(), childDataType.Fmt()))
	}

	return boolDataType, nil
}

//getAsterisk is used to inform the data type of expressions with
//"*" as operators of higher precedence as well as to validate that its sub-expressions are of a valid data type.
func (analyzer *SemanticAnalyzer) getAsterisk(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	numberChildren := len(tree.Current.Children)
	switch numberChildren{
	case 1: return analyzer.getByCategory(tree, scope)
	case 2: return analyzer.getOperatorLeftGreaterThanRightSize(tree, scope)
	default:
		//Due to the syntactic analyzer, the tree should not allow the current node to have anything but one or two children
		panic(errorhandler.UnexpectedCompilerError())
	}

}

//getArray is used to inform the data type of an arrays well as to validate that its sub-expressions are of a valid data type.
func (analyzer *SemanticAnalyzer) getArray(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	return analyzer.getByCategory(tree, scope)
}

//getAddress is used to inform the data type of a memory address as well as to validate that its sub-expressions are of a valid data type.
func (analyzer *SemanticAnalyzer) getAddress(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	//Due to the syntax analyzer, the tree should not allow the current node to have anything but one child.
	if len(tree.Current.Children) != 1{
		panic(errorhandler.UnexpectedCompilerError())
	}
	tree.Current = tree.Current.Children[0]
	pointsTo, err := analyzer.getIdent(tree, scope)
	if err != nil{
		return nil, err
	}

	return token.NewDataType(token.DataTypePointer, 2, 0, pointsTo), nil
}

//FindLeaf is called when searching an identifier or a keyword of a data type after a sequenve of arrays and pointers
// In order to be in this method, the node of the tree that is being analyzed must be an identifier, a "]" or a "*".
func FindLeaf(tree *ast.SyntaxTree) *ast.Node{
	//Due to the syntax analyzer and the context of this function, the tree should not allow the current node to have anything but zero, one or two children.
	if len(tree.Current.Children)  == 0{
		return tree.Current
	}
	if len(tree.Current.Children)  > 2{
		panic(errorhandler.UnexpectedCompilerError())
	}


	//to find an identifier or a keyword of a data type after a sequence of arrays and pointers, we must find the leaf
	//Due to the syntactic analyzer the leaf can be found searching through the only child of a "*" and the right child of a "]"
	leaf := tree.Current.Children[len(tree.Current.Children)-1]
	for len(leaf.Children) != 0{
		leaf = leaf.Children[len(leaf.Children)-1]
	}
	return leaf

}

//getByCategory analyzes the context of a data type (whether it is the data type of an identifier
//that should already be saved in the symbol table, or the data type in a declaration)
//and redirects the analysis to the corresponding method.Then it returns the data type returned by those methods.
// In order to be in this method, the node of the tree that is being analyzed must be a "]" or a "*".
func (analyzer *SemanticAnalyzer) getByCategory(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	leaf := FindLeaf(tree)
	if leaf.Value.Type == token.IDENT {
		return analyzer.getIdent(tree, scope)

	}
	if leaf.Value.Type == token.TYPEBOOL ||  leaf.Value.Type == token.TYPEBYTE || leaf.Value.Type == token.VOID {
		return analyzer.getNewIdent(tree, scope)


	}

	//In the context of this method a leaf should be always a identifier or
	//a keyword to refer a data type declaration.
	panic(errorhandler.UnexpectedCompilerError())

}

//getIdent checks whether we are analysing a reference to the symbol table or a deference to a identifier,
//then it redirects the analysis to the corresponding method. Then it returns the data type returned by those methods
func (analyzer *SemanticAnalyzer) getIdent(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	identDataType, err := analyzer.getReferencedIdent(tree, scope)
	if err != nil{
		return nil, err
	}

	//if tree.Current is not a leaf, then we are analysing a name variable and we not need a further analysis of the data type
	if len(tree.Current.Children) == 0{
		return identDataType, nil
	}else{
		//if instead we are analysing a temporary value that deference a variable stored in the symbol table, we keep analysing
		return analyzer.getDereferencedIdent(tree, scope, identDataType)
	}

}

//getReferencedIdent checks if a identifier is in the symbol table, validates if it's or not a function,
//and returns the data type of the identifier.
func (analyzer *SemanticAnalyzer) getReferencedIdent(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	leaf := FindLeaf(tree)
	symbol, ok := analyzer.symbolTable[scope+leaf.Value.Literal]
	if !ok {
		return nil, errors.New(errorhandler.UnresolvedReference(leaf.Value.Line, leaf.Value.Literal))
	}

	analyzer.isValidByFunctionContext(tree, symbol)

	if !analyzer.isValidByFunctionContext(tree, symbol) {
		if symbol.IsFunction {
			return nil, errors.New(errorhandler.IdentifierIsFunction(leaf.Value.Line, leaf.Value.Literal))
		} else {
			return nil, errors.New(errorhandler.IdentifierIsNotFunction(leaf.Value.Line, leaf.Value.Literal))

		}
	}
	return &symbol.DataType, nil
}

//getDereferencedIdent checks if a temporary value, which is the result of dereferencing an identifier,
//has a valid data type and returns it.
//For example, if "foo" is a variable of type "pointer to pointer to array", and we are analyzing "*foo",
//we recognize it as a valid data type, which is "pointer to array".
//If, instead, we analyze "[i]foo", we recognize that it is not a valid data type
//(because dereferences of identifiers can only be read from left to right).
func (analyzer *SemanticAnalyzer) getDereferencedIdent(tree *ast.SyntaxTree, scope string,
	identifierDataType *token.DataType) (*token.DataType, error) {
	dataTypeToCompare := identifierDataType
	nodeToAnalyse := tree.Current
	identifier := FindLeaf(tree)
	for nodeToAnalyse != identifier {
		//we skip parentheses
		if nodeToAnalyse.Value.Type == token.RPAREN{
			nodeToAnalyse = nodeToAnalyse.Children[0]
			continue
		}
		switch dataTypeToCompare.Kind {
		case token.DataTypePointer:
			if nodeToAnalyse.Value.Type == token.ASTERISK{
				dataTypeToCompare = dataTypeToCompare.PointsTo
				nodeToAnalyse = nodeToAnalyse.Children[0]
				continue
			}else{
				return nil, errors.New(errorhandler.InvalidIndirectOf(nodeToAnalyse.Value.Line, identifier.Value.Literal))
			}
		case token.DataTypeArray:
			if nodeToAnalyse.Value.Type == token.RBRACKET {

				err := analyzer.checkIndexArray(tree, scope, dataTypeToCompare.Length, nodeToAnalyse)
				if err != nil {
					return nil, err
				}
			}else {
				return nil, errors.New(errorhandler.InvalidIndirectOf(nodeToAnalyse.Value.Line, identifier.Value.Literal))

			}
			dataTypeToCompare = dataTypeToCompare.PointsTo
			nodeToAnalyse = nodeToAnalyse.Children[1]
			continue
		default:
			//In the context of this method, we only consider a sequence of arrays and pointers before an identifier to be valid.
			return nil, errors.New(errorhandler.InvalidIndirectOf(nodeToAnalyse.Value.Line, identifier.Value.Literal))

		}
	}
	return dataTypeToCompare, nil
}

//TODO: VEEEEER
//checkIndexArray checks if the index of an array is valid, meaning that it needs to be a byte, and in the case of literals it cannot be out of bounds
//when the context allows this error.
func (analyzer *SemanticAnalyzer) checkIndexArray(tree *ast.SyntaxTree, scope string,
	maxLength int, bracket *ast.Node) error {

	//Due to the syntax analyzer, the tree should not allow the current node to have anything but two children.
	if len(bracket.Children) != 2 {
		panic(errorhandler.UnexpectedCompilerError())
	}
	if bracket.Children[0].Value.Type == token.BYTE {
		index, err := strconv.Atoi(bracket.Children[0].Value.Literal)
		if err != nil {
			return  errors.New(errorhandler.IndexMustBeAByte(bracket.Value.Line))
		} else {
			if maxLength != 0{
				if index >= maxLength {
					return  errors.New(errorhandler.IndexOutOfBounds(bracket.Value.Line))
				}

			}
		}
	} else {
		current := tree.Current
		tree.Current = tree.Current.Children[0]
		indexDataType, err := analyzer.getDataType[bracket.Children[0].Value.Type](tree, scope)
		tree.Current = current
		if err != nil {
			return  err
		} else {
			if indexDataType != token.NewDataType(token.DataTypeByte, 1, 0, nil) {
				return  errors.New(errorhandler.IndexMustBeAByte(bracket.Children[0].Value.Line))

			}
		}
	}


	return nil
}

//isValidByFunctionContext checks the context of an identifier.
// It returns "true" if the identifier is in a function call and represents a function, or if it is not in a function call and represents a variable.
//It returns "false" otherwise.
func (analyzer *SemanticAnalyzer) isValidByFunctionContext(tree *ast.SyntaxTree, symbol token.Symbol) bool{

	//This function only can be used in the context of a variable or a dereference to a variable
	if tree.Current.Value.Type == token.DOLLAR || tree.Current.Value.Type == token.ASTERISK ||
		tree.Current.Value.Type == token.RBRACKET || tree.Current.Value.Type == token.IDENT{
		if tree.Current.Parent.Value.Type == token.RPAREN{
			if symbol.IsFunction {
				return true
			}else{
				return false
			}
		}else{
			if symbol.IsFunction {
				return false
			}else{
				return true
			}
		}

	}
	panic(errorhandler.UnexpectedCompilerError())

}

//getNewIdent informs the data type of an identifier that is being declared
func (analyzer *SemanticAnalyzer) getNewIdent(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	current := tree.Current
	switch current.Value.Type {
	case token.ASTERISK:
		if len(current.Children) != 1 {
			panic(errorhandler.UnexpectedCompilerError())
		}
		tree.Current = current.Children[0]
		pointsTo, err := analyzer.getNewIdent(tree, scope)
		tree.Current = current
		if err != nil {
			return nil, err
		}

		return token.NewDataType(token.DataTypePointer, 2, 0, pointsTo), nil
	case token.RBRACKET:
		if len(current.Children) != 2 {
			panic(errorhandler.UnexpectedCompilerError())
		}
		length := 0

		err := analyzer.checkIndexArray(tree, scope,0, tree.Current)
		tree.Current = current.Children[1]
		pointsTo, err := analyzer.getNewIdent(tree, scope)
		if err != nil {
			return nil, err
		}
		tree.Current = current
		return token.NewDataType(token.DataTypeArray, 2, length, pointsTo), nil

	case token.TYPEBYTE:
		return analyzer.getByte(tree, scope)
	case token.TYPEBOOL:
		return analyzer.getBool(tree, scope)
	case token.VOID:
		return analyzer.getVoid(tree, scope)
	default:
		panic(errorhandler.UnexpectedCompilerError())

	}
}

//getByte is used to inform the data type of bytes
func (analyzer *SemanticAnalyzer) getByte(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	return token.NewDataType(token.DataTypeByte, 1, 0, nil), nil
}

//getBool is used to inform the data type of booleans
func (analyzer *SemanticAnalyzer) getBool(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	return token.NewDataType(token.DataTypeBool, 1, 0, nil), nil
}


//getVoid is used to inform the data type of booleans
func (analyzer *SemanticAnalyzer) getVoid(tree *ast.SyntaxTree, scope string) (*token.DataType, error) {
	return token.NewDataType(token.DataTypeVoid, 1, 0, nil), nil
}