package errorhandler

import (
	"strconv"
)

func UnexpectedDataType(line int, expected string, unexpected string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nExpected: " + expected + " got: " + unexpected
	return errorString
}

func PointerToVoid(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nPointer to void. "
	return errorString

}
func DataTypesMismatch(line int, actualDatatype string, symbol string, expectedDatatype string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nData types mismatches: " + actualDatatype + " " + symbol + " " + expectedDatatype
	return errorString

}
func ByteOutOfRange(line int, number int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nNumber: " + strconv.Itoa(number) + " is not a byte"
	return errorString

}

func InvalidAssignation(line int, datatype string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nInvalid assignation to: " + datatype
	return errorString

}

func UnexpectedCompilerError() string {
	errorString := "\nunexpected compiler error\n"
	return errorString
}
func NumberOfParametersDoesntMatch(line int, actualLength int, expectedLength int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nThe number of parameters doesn't match " +
		strconv.Itoa(actualLength) + "=" + strconv.Itoa(expectedLength) + "\n"
	return errorString
}

func UnresolvedReference(line int, reference string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nUnresolved reference: " + reference
	return errorString
}

func UnallowedPointerToArray(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nUnallowed pointer to array"

	return errorString
}
func InvalidIndirectOf(line int, reference string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nInvalid indirect of: " + reference
	return errorString
}

func IndexOutOfBounds(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nIndex out of bound"
	return errorString
}

func IndexMustBeAByte(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nThe index of an array must be a byte"
	return errorString
}

func IdentifierIsFunction(line int, reference string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nIdentifier " + reference + " is a function"
	return errorString
}

func IdentifierIsNotFunction(line int, reference string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nIdentifier " + reference + " is not a function"
	return errorString
}

func IdentifierMissed(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nIdentifier missed"
	return errorString
}
func NameAlreadyInUse(line int, reference string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nThe name " + reference + " is already in use"
	return errorString
}

func NegativeIndex(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nNegative index."
	return errorString
}

func UnreachableCode(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) +
		"\nUnreachable code "
	return errorString

}
func IllegalToken(line int, t string) string {

	errorString := "\n illegal token: \"" + t + "\" \n in line: " + strconv.Itoa(line)
	return errorString
}

func FunctionOutsideGlobalScope(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) + "\nfunction declaration outside global scope"
	return errorString
}

func GlobalScopeOnlyAllowsDeclarations(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) + "\nglobal scope only allows declarations"
	return errorString
}
func MainFunctionNeeded() string {
	errorString := "main function needed\n"
	return errorString
}
func SyntaxError() string {
	errorString := "syntactic error"
	return errorString
}

func InvalidReturnType(line int, returnType string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) + "\n: " + returnType + " is a invalid return type"
	return errorString
}

func InvalidParamType(line int, paramType string) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) + "\n: " + paramType + " is a invalid parameter"
	return errorString
}

func TooManyParams(line int) string {
	errorString := "semantic error\nin line: " + strconv.Itoa(line) + "\nparams exceed the limit of nine bytes"
	return errorString
}

func TooManyRegisters(line int) string {
	errorString := "error\nin line: " + strconv.Itoa(line) + "\nthe expression requires too many registers to be solve"
	return errorString
}

func NotEnoughMemory() string {
	errorString := "Not Enough Memory"
	return errorString
}
