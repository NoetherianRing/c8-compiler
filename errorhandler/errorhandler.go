package errorhandler

import (
	"strconv"
)

func UnexpectedDataType(line int, expected string, unexpected string) string{
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nExpected: " + expected + " got: " + unexpected
	return errorString
}

func PointerToVoid(line int) string{
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nPointer to void. "
	return errorString

}
func DataTypesMismatch(line int, actualDatatype string, symbol string, expectedDatatype string) string{
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nData types mismatches: " + actualDatatype + " " + symbol + " " + expectedDatatype
	return errorString

}
func UnexpectedCompilerError() string {
	errorString := "\nunexpected compiler error\n"
	return errorString
}
func NumberOfParametersDoesntMatch(line int, actualLength int, expectedLength int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nThe number of parameters doesn't match " 	+
		strconv.Itoa(actualLength) + "=" +  strconv.Itoa(expectedLength) + "\n"
	return errorString
}

func UnresolvedReference(line int, reference string) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nUnresolved reference: " + reference
	return errorString
}

func InvalidIndirectOf(line int, reference string) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nInvalid indirect of: " + reference
	return errorString
}

func IndexOutOfBounds(line int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nIndex out of bound"
	return errorString
}

func IndexMustBeAByte(line int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nThe index of an array must be a byte"
	return errorString
}

func IdentifierIsFunction(line int, reference string) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nIdentifier " + reference + " is a function"
	return errorString
}


func IdentifierIsNotFunction(line int, reference string) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nIdentifier " + reference + " is not a function"
	return errorString
}

func IdentifierMissed(line int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nIdentifier missed"
	return errorString
}
func NameAlreadyInUse(line int, reference string) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nThe name " + reference + " is already in use"
	return errorString
}

func NegativeIndex(line int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nNegative index."
	return errorString
}

func UnreachableCode(line int) string{
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nUnreachable code "
	return errorString

}