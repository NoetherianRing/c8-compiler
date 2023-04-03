package semanticAnalyzer

import (
	"errors"
	"github.com/NoetherianRing/c8-compiler/ast"
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"github.com/NoetherianRing/c8-compiler/symboltable"
	"github.com/NoetherianRing/c8-compiler/token"
	"strconv"
)

type DataTypeFactory struct {
	scope        *symboltable.Scope
	ctxNode      *ast.Node
	walkingAFunc bool
}

func NewDataTypeFactory() *DataTypeFactory {
	getter := new(DataTypeFactory)
	getter.walkingAFunc = false
	getter.ctxNode = nil
	getter.scope = nil
	return getter
}

func (getter *DataTypeFactory) SetCxtNode(node *ast.Node) {
	getter.ctxNode = node
}

func (getter *DataTypeFactory) SetScope(scope *symboltable.Scope) {
	getter.scope = scope
}

//GetDataType calls "redirect()"to obtain a function that returns the data type of the current node of the tree.
//It then executes that function and returns its result.
func (getter *DataTypeFactory) GetDataType() (interface{}, error) {
	if getter.ctxNode == nil || getter.scope == nil {
		panic(errorhandler.UnexpectedCompilerError())
	}
	if getter.scope == nil {
		panic(errorhandler.UnexpectedCompilerError())
	}
	get := getter.redirect()
	return get()
}

//redirect analyzes the token type of the current node of the tree and returns a function that analyzes the data type of the expression led by the node.
//When needed, it calls another redirect method for further analysis.
func (getter *DataTypeFactory) redirect() func() (interface{}, error) {

	switch getter.ctxNode.Value.Type {
	case token.ASTERISK:
		return getter.redirectAsterisk()
	case token.RBRACKET:
		return getter.redirectBracket()
	case token.RPAREN:
		return getter.redirectParentheses()
	case token.VOID:
		return getter.declarationSimple
	case token.TYPEBYTE:
		return getter.declarationSimple
	case token.TYPEBOOL:
		return getter.declarationSimple
	case token.IDENT:
		return getter.reference
	case token.LOR:
		return getter.logicExpression
	case token.LAND:
		return getter.logicExpression
	case token.BANG:
		return getter.logicExpression
	case token.AND:
		return getter.bitwiseExpression
	case token.XOR:
		return getter.bitwiseExpression
	case token.OR:
		return getter.bitwiseExpression
	case token.PLUS:
		return getter.numericOperation
	case token.MINUS:
		return getter.numericOperation
	case token.LTLT:
		return getter.byteOperation
	case token.GTGT:
		return getter.byteOperation
	case token.SLASH:
		return getter.byteOperation
	case token.PERCENT:
		return getter.byteOperation
	case token.EQEQ:
		return getter.comparison
	case token.NOTEQ:
		return getter.comparison
	case token.LT:
		return getter.numericLogicalComparison
	case token.LTEQ:
		return getter.numericLogicalComparison
	case token.GT:
		return getter.numericLogicalComparison
	case token.GTEQ:
		return getter.numericLogicalComparison
	case token.DOLLAR:
		return getter.address
	case token.BYTE:
		return getter.simple
	case token.BOOL:
		return getter.simple
	default:
		panic(errorhandler.UnexpectedCompilerError())
	}

}

//redirectAsterisk returns a function that analysis the data type of an expression led by an asterisk by checking the context.
func (getter *DataTypeFactory) redirectAsterisk() func() (interface{}, error) {

	switch len(getter.ctxNode.Children) {
	case 1:
		if getter.isADeclarationContext() {
			return getter.declaration
		} else {
			return getter.dereference
		}
	case 2:
		return getter.byteOperation

	default:
		panic(errorhandler.UnexpectedCompilerError())
	}

}

//redirectBracket returns a function that analysis the data type of an expression led by an bracket by checking the context.
func (getter *DataTypeFactory) redirectBracket() func() (interface{}, error) {
	if getter.isADeclarationContext() {
		return getter.declaration
	} else {
		return getter.dereference
	}

}

//redirectParentheses returns a function that analysis the data type of an expression led by a parentheses by checking the context.
func (getter *DataTypeFactory) redirectParentheses() func() (interface{}, error) {
	child := getter.ctxNode.Children[0].Value
	if child.Type == token.IDENT {
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

	if leafType == token.TYPEBOOL || leafType == token.TYPEBYTE {
		return true
	}
	return false
}

//validateParamsDataType validates if the data type of the parameters of a function call match with each data type
//of expressions led by the ctxNode. It also validates if the numbers of parameters matches.
func (getter *DataTypeFactory) validateParamsDataType(args []interface{}) error {
	var err error
	err = nil
	moreParams := true
	backup := getter.ctxNode
	for i := range args {
		if getter.ctxNode.Value.Type == token.COMMA {

			comma := getter.ctxNode
			getter.ctxNode = getter.ctxNode.Children[0]
			err = getter.validateParamDataType(args, i)
			if err != nil {
				return err
			}

			getter.ctxNode = comma.Children[1]

		} else {
			moreParams = false
			if len(args)-(i) != 1 {
				line := getter.ctxNode.Value.Line
				err = errors.New(errorhandler.NumberOfParametersDoesntMatch(line, i+1, len(args)))
				return err
			} else {
				err = getter.validateParamDataType(args, i)
			}
		}

	}
	if moreParams {
		getter.ctxNode = backup
		numberOfTreeParams := CountNodesToLeafByRight(getter.ctxNode) + 1
		line := getter.ctxNode.Value.Line
		err = errors.New(errorhandler.NumberOfParametersDoesntMatch(line, numberOfTreeParams, len(args)))

	}

	return err
}

//validateParamDataType  validates if the data type of the expression led by the current "ctxNode"
//matches the data type of the argument "i" of a function.
func (getter *DataTypeFactory) validateParamDataType(args []interface{}, i int) error {
	treeParam, err := getter.GetDataType()
	if err != nil {
		return err
	}
	if symboltable.Compare(treeParam, args[i]) {
		return nil
	} else {

		line := getter.ctxNode.Value.Line
		err = errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(treeParam), token.EQ, symboltable.Fmt(args[i])))
		return err
	}

}

//logicExpression verifies that the expressions led by the ctx Node are boolean and returns a error if not.
//Otherwise returns a boolean
func (getter *DataTypeFactory) logicExpression() (interface{}, error) {
	boolType := symboltable.NewBool()
	for _, child := range getter.ctxNode.Children {
		backup := getter.ctxNode
		getter.ctxNode = child
		childDatatype, err := getter.GetDataType()
		getter.ctxNode = backup
		if err != nil {
			return nil, err
		}
		if !boolType.Compare(childDatatype) {
			line := child.Value.Line
			err := errors.New(errorhandler.UnexpectedDataType(line, symboltable.Fmt(boolType), symboltable.Fmt(childDatatype)))
			return nil, err
		}
	}
	return boolType, nil
}

//comparison verifies that the expressions led by the ctx Node are of the same data type and returns a error if not.
//Otherwise returns a boolean
func (getter *DataTypeFactory) comparison() (interface{}, error) {
	backup := getter.ctxNode
	getter.ctxNode = getter.ctxNode.Children[0]
	leftChildDataType, err := getter.GetDataType()
	getter.ctxNode = backup
	if err != nil {
		return nil, err
	}
	getter.ctxNode = getter.ctxNode.Children[1]
	rightChildDataType, err := getter.GetDataType()
	getter.ctxNode = backup
	if err != nil {
		return nil, err
	}
	if leftChildDataType != rightChildDataType {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(leftChildDataType), token.EQ, symboltable.Fmt(rightChildDataType)))
		return nil, err
	}
	return symboltable.NewBool(), nil

}

//numericLogicalComparison verifies that the expressions led by the ctx Node are of the same numeric data type by calling to
//validateSameNumericDataType(). Returns a error if needed, otherwise returns a boolean
func (getter *DataTypeFactory) numericLogicalComparison() (interface{}, error) {
	_, err := getter.validateSameNumericDataType()
	if err != nil {
		return nil, err
	}
	return symboltable.NewBool(), nil
}

//numericLogicalComparison verifies that the expressions led by the ctx Node are of the same numeric data type by calling to
//validateSameNumericDataType(). Returns a error if needed, otherwise returns a the same datatype of the expressions it leads.
func (getter *DataTypeFactory) bitwiseExpression() (interface{}, error) {
	datatype, err := getter.validateSameNumericDataType()
	if err != nil {
		return nil, err
	}
	return datatype, nil
}

//validateSameNumericDataType verifies that the expressions led by the ctx Node are of the same numeric data type and returns a error if not.
//Otherwise returns the same datatype of the expressions it leads.
func (getter *DataTypeFactory) validateSameNumericDataType() (interface{}, error) {

	leftChildDataType, rightChildDataType, err := getter.obtainOperandsDatatype()
	if err != nil {
		return nil, err
	}
	if !symboltable.Compare(leftChildDataType, rightChildDataType) {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.DataTypesMismatch(line, symboltable.Fmt(leftChildDataType), token.EQ, symboltable.Fmt(rightChildDataType)))
		return nil, err
	} else {
		if !symboltable.IsNumeric(leftChildDataType) {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.UnexpectedDataType(line, "numeric", symboltable.Fmt(leftChildDataType)))
			return nil, err
		}

	}

	return leftChildDataType, nil
}

//obtainOperandDatatype return the datatype of two operands of a operation and an error if needed
func (getter *DataTypeFactory) obtainOperandsDatatype() (interface{}, interface{}, error) {
	backup := getter.ctxNode
	getter.ctxNode = getter.ctxNode.Children[0]
	leftChildDataType, err := getter.GetDataType()
	if err != nil {

		return nil, nil, err
	}
	getter.ctxNode = backup
	getter.ctxNode = getter.ctxNode.Children[1]

	rightChildDataType, err := getter.GetDataType()
	if err != nil {
		return nil, nil, err
	}
	return leftChildDataType, rightChildDataType, nil
}

//numericOperation verifies that the left child of ctx Node is a pointer or a byte and the right child a byte
func (getter *DataTypeFactory) numericOperation() (interface{}, error) {
	leftChildDataType, rightChildDataType, err := getter.obtainOperandsDatatype()
	if err != nil {
		return nil, err
	}

	if !symboltable.IsNumeric(leftChildDataType) {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnexpectedDataType(line, "numeric", symboltable.Fmt(leftChildDataType)))
		return nil, err
	}
	if !symboltable.IsByte(rightChildDataType) {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnexpectedDataType(line, "byte", symboltable.Fmt(rightChildDataType)))
		return nil, err

	}

	return leftChildDataType, nil
}

//byteOperation verifies that the left  and right children of ctx Node are both bytes
func (getter *DataTypeFactory) byteOperation() (interface{}, error) {
	leftChildDataType, rightChildDataType, err := getter.obtainOperandsDatatype()
	if err != nil {
		return nil, err
	}

	if !symboltable.IsByte(leftChildDataType) {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnexpectedDataType(line, "byte", symboltable.Fmt(leftChildDataType)))
		return nil, err
	}
	if !symboltable.IsByte(rightChildDataType) {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnexpectedDataType(line, "byte", symboltable.Fmt(rightChildDataType)))
		return nil, err

	}

	return leftChildDataType, nil
}

//address takes the data type of what ctxNode dereference and return a pointer that points to
//that data type
func (getter *DataTypeFactory) address() (interface{}, error) {
	getter.ctxNode = getter.ctxNode.Children[0]
	pointsTo, err := getter.dereference()
	if err != nil {
		return nil, err
	}
	return symboltable.NewPointer(pointsTo), nil
}

//functionCall validates that the data types of the parameters of a function call
//match the ones of its definition and returns an error if needed.
//If not, it returns the data type of the return value of the function.
func (getter *DataTypeFactory) functionCall() (interface{}, error) {
	getter.walkingAFunc = true

	backup := getter.ctxNode
	identifier := getter.ctxNode.Children[0]

	getter.ctxNode = identifier
	identifierDataType, err := getter.reference()
	getter.ctxNode = backup

	getter.walkingAFunc = false

	if err != nil {
		return nil, err
	}

	funcDataType, ok := identifierDataType.(symboltable.Function)
	if !ok {
		panic(errorhandler.UnexpectedCompilerError())
	}

	returnDataType := funcDataType.Return
	argsDataType := funcDataType.Args

	if argsDataType != nil {
		if len(getter.ctxNode.Children) != 2 {
			line := getter.ctxNode.Value.Line
			err = errors.New(errorhandler.NumberOfParametersDoesntMatch(line, 0, len(argsDataType)))
			return nil, err
		}

		param := getter.ctxNode.Children[1]
		getter.ctxNode = param
		err = getter.validateParamsDataType(argsDataType)
		getter.ctxNode = backup
		if err != nil {
			return nil, err
		}
	} else {
		if len(getter.ctxNode.Children) > 1 {
			getter.ctxNode = getter.ctxNode.Children[1]
			numberOfTreeParams := CountNodesToLeafByRight(getter.ctxNode) + 1
			line := getter.ctxNode.Value.Line
			err = errors.New(errorhandler.NumberOfParametersDoesntMatch(line, numberOfTreeParams, len(argsDataType)))
			return nil, err
		}
	}
	return returnDataType, nil

}

//reference checks if an identifier is stored in the symbol table and, if it is, returns its data type.
//It returns an error in case of not finding the reference, if we are expecting a function and the reference is not a function,
//or if we are not expecting a function and the reference is a function.
func (getter *DataTypeFactory) reference() (interface{}, error) {
	literal := getter.ctxNode.Value.Literal
	ref, ok := getter.scope.Symbols[literal]
	if !ok {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnresolvedReference(line, literal))
		return nil, err
	} else {
		if getter.walkingAFunc && !ref.IsFunction {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.IdentifierIsNotFunction(line, literal))
			return nil, err
		}

		if !getter.walkingAFunc && ref.IsFunction {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.IdentifierIsFunction(line, literal))
			return nil, err
		}
	}
	return ref.DataType, nil
}

//dereference analyzes the data type of a dereference and returns it
func (getter *DataTypeFactory) dereference() (interface{}, error) {
	identifier := GetLeafByRight(getter.ctxNode)
	backup := getter.ctxNode
	if identifier.Value.Type != token.IDENT {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.IdentifierMissed(line))
		return nil, err
	}
	if identifier.Parent != nil {
		if identifier.Parent.Value.Type == token.RPAREN {
			getter.ctxNode = identifier.Parent
		} else {
			getter.ctxNode = identifier
		}
	}
	toCompare, err := getter.GetDataType()
	getter.ctxNode = backup
	if err != nil {
		return nil, err
	}

	for getter.ctxNode != identifier {

		switch toCompare.(type) {
		case symboltable.Pointer:
			if getter.ctxNode.Value.Type == token.ASTERISK {
				getter.ctxNode = getter.ctxNode.Children[0]
				if getter.ctxNode.Value.Type == token.RPAREN {
					getter.ctxNode = getter.ctxNode.Children[0]

				}
				toCompare = toCompare.(symboltable.Pointer).PointsTo
			} else {
				line := getter.ctxNode.Value.Line
				err := errors.New(errorhandler.InvalidIndirectOf(line, identifier.Value.Literal))
				return nil, err

			}
		case symboltable.Array:
			if getter.ctxNode.Value.Type == token.RBRACKET {
				backup = getter.ctxNode
				getter.ctxNode = getter.ctxNode.Children[0]
				err := getter.validateIndex(toCompare)
				getter.ctxNode = backup
				if err != nil {
					return nil, err
				}
				getter.ctxNode = getter.ctxNode.Children[1]
				if getter.ctxNode.Value.Type == token.RPAREN {
					getter.ctxNode = getter.ctxNode.Children[0]

				}
				toCompare = toCompare.(symboltable.Array).Of

			} else {
				line := getter.ctxNode.Value.Line
				err := errors.New(errorhandler.InvalidIndirectOf(line, identifier.Value.Literal))
				return nil, err
			}
		default:
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.InvalidIndirectOf(line, identifier.Value.Literal))
			return nil, err
		}

	}

	return toCompare, nil
}

//validateIndex validates if the index of an array is a byte and if its out of bound.
func (getter *DataTypeFactory) validateIndex(compare interface{}) error {

	if getter.ctxNode.Value.Type != token.BYTE {
		return errors.New(errorhandler.UnexpectedCompilerError())

	} else {
		length, err := strconv.Atoi(getter.ctxNode.Value.Literal)
		if err != nil {
			panic(errorhandler.UnexpectedCompilerError())
		}
		if length < 0 {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.NegativeIndex(line))
			return err
		}
		arrayToCompare := compare.(symboltable.Array)
		if length >= arrayToCompare.Length {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.IndexOutOfBounds(line))
			return err
		}
	}

	return nil
}

//declaration verifies that there is no declaration of a pointer to a function
//in the context and returns an error if needed. Otherwise, it returns the data type built by "declarationFactory()".
func (getter *DataTypeFactory) declaration() (interface{}, error) {
	if len(getter.ctxNode.Children) != 0 {
		leaf := GetLeafByRight(getter.ctxNode)
		leafType := leaf.Value.Type
		if leafType == token.VOID {
			line := getter.ctxNode.Value.Line
			err := errors.New(errorhandler.PointerToVoid(line))
			return nil, err
		}
	}
	return getter.declarationFactory()
}

//declarationFactory builds the datatype in the context of a declaration by calling more specific methods
func (getter *DataTypeFactory) declarationFactory() (interface{}, error) {
	if len(getter.ctxNode.Children) == 0 {
		return getter.declarationSimple()
	} else {
		switch getter.ctxNode.Value.Type {
		case token.ASTERISK:
			return getter.declarationFactoryPointer()
		case token.RBRACKET:
			return getter.declarationFactoryArray()
		default:
			panic(errorhandler.UnexpectedCompilerError())

		}
	}
}

//declarationSimple returns a boolean, a byte or a void depending on the context
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

//declarationFactoryPointer returns a pointer data type. The data type of the reference it points to, is given by moving
//the context to the child of the current node and building its datatype
func (getter *DataTypeFactory) declarationFactoryPointer() (interface{}, error) {
	getter.ctxNode = getter.ctxNode.Children[0]
	if getter.ctxNode.Value.Type == token.RBRACKET {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.UnallowedPointerToArray(line))
		return nil, err
	} else {
		pointsTo, err := getter.declarationFactory()
		if err != nil {
			return nil, err
		}
		return symboltable.NewPointer(pointsTo), nil

	}
}

//declarationFactoryArray validates that the index of an array is valid (by checking its data type) and return a array data type
//the data type of the elements of the array is obtained by moving the context and calling to declarationFactory()
func (getter *DataTypeFactory) declarationFactoryArray() (interface{}, error) {
	index := getter.ctxNode.Children[0]
	if index.Value.Type != token.BYTE {
		return nil, errors.New(errorhandler.UnexpectedCompilerError())
	}

	literal, err := strconv.Atoi(index.Value.Literal)
	if err != nil {
		panic(errorhandler.UnexpectedCompilerError())
	}
	length := literal
	if length < 0 {
		line := getter.ctxNode.Value.Line
		err := errors.New(errorhandler.NegativeIndex(line))
		return nil, err
	}

	getter.ctxNode = getter.ctxNode.Children[1]
	of, err := getter.declarationFactory()
	if err != nil {
		return nil, err
	}
	return symboltable.NewArray(length, of), nil
}

//simple returns a boolean or a byte depending on the context
func (getter *DataTypeFactory) simple() (interface{}, error) {
	switch getter.ctxNode.Value.Type {
	case token.BOOL:
		return symboltable.NewBool(), nil
	case token.BYTE:
		_byte, err := strconv.Atoi(getter.ctxNode.Value.Literal)
		if err != nil {
			return nil, errors.New(errorhandler.UnexpectedCompilerError())
		}
		if _byte > 255 || _byte < 0 {
			line := getter.ctxNode.Value.Line
			return nil, errors.New(errorhandler.ByteOutOfRange(line, _byte))

		}
		return symboltable.NewByte(), nil
	default:
		panic(errorhandler.UnexpectedCompilerError())

	}
}

// GetLeafByRight gets the leaf by walking a tree using the right child of each node.
func GetLeafByRight(head *ast.Node) *ast.Node {
	current := head
	for len(current.Children) != 0 {
		current = current.Children[len(current.Children)-1]
	}
	return current
}

// CountNodesToLeafByRight counts how many nodes are left to found
//a leaf by walking a tree using the right child of each node.
func CountNodesToLeafByRight(head *ast.Node) int {
	current := head
	i := 0
	for len(current.Children) != 0 {
		i++
		current = current.Children[len(current.Children)-1]
	}
	return i
}
