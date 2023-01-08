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
func DataTypesMismatch(line int, dataType1 string, symbol string, dataType2 string) string{
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nData types mismatches: " + dataType1 + " " + symbol + " " +dataType2
	return errorString

}
func UnexpectedCompilerError() string {
	errorString := "\nunexpected compiler error\n"
	return errorString
}
func NumberOfParametersDoesntMatch(line int, len1 int, len2 int) string {
	errorString := "semantic error\nin line: "+ strconv.Itoa(line) +
		"\nThe number of parameters doesn't match " 	+
		strconv.Itoa(len1) + "=" +  strconv.Itoa(len2) + "\n"
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