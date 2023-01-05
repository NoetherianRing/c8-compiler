package token

import "strconv"

type SymbolTable map[string] Symbol

const (
	DataTypeByte = iota
	DataTypeBool
	DataTypeArray
	DataTypePointer
	DataTypeVoid
)

type DataType struct {
	Kind     int
	Size     int
	Length   int
	PointsTo *DataType
}

func (dt DataType) Fmt() string{
	switch dt.Kind {
	case DataTypeByte:
		return "byte"
	case DataTypeBool:
		return "bool"
	case DataTypeArray:
		return "["+ strconv.Itoa(dt.Length)+"]" + dt.PointsTo.Fmt()
	case DataTypePointer:
		return "*"+dt.PointsTo.Fmt()
	case DataTypeVoid:
		return "void"
	default:
		return "illegal"
	}
}

func NewDataType(kind int, size int, length int, pointsTo *DataType) *DataType{
	return &DataType{Kind: kind,
		Size: size,
		Length: length,
		PointsTo: pointsTo,
	}
}

type Symbol struct{
	Symbol     string
	DataType   DataType
	Scope      string
	IsFunction bool
	Args       []DataType

}

func NewSymbol(symbol string, datatype DataType, scope string, isFunction bool, args []DataType) Symbol {
	return Symbol{
		Symbol:     symbol,
		DataType:   datatype,
		Scope:      scope,
		IsFunction: isFunction,
		Args:       args,
	}
}

