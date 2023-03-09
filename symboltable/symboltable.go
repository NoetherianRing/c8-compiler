package symboltable

import (
	"github.com/NoetherianRing/c8-compiler/errorhandler"
	"reflect"
	"strconv"
)

type SymbolTable map[string]*Symbol

const(
	KindByte = iota
	KindVoid
	KindBool
)

type Scope struct{
	SubScopes []*Scope
	NumberOfSubScope int
	Parent *Scope
	Symbols SymbolTable
}

type Simple struct{
	Size int
	Kind int
}
type Array struct{
	Length int
	Of interface{}
}
type Pointer struct{
	Size int
	PointsTo interface{}
}

type Function struct{
	Return interface{}
	Args [] interface{}
}

type Symbol struct{
	Identifier string
	IsFunction bool
	DataType interface{}
}
func (array Array) SizeOfElements() int{
	switch array.Of.(type) {
	case Pointer:
		return array.Of.(Pointer).Size
	case Simple:
		return array.Of.(Simple).Size
	case Array:
		of := array.Of.(Array)
		return of.SizeOfElements() * of.Length
	default:
		return 0
	}

}

func (t Simple)Compare(datatype interface{})bool{
	toCompare, ok := datatype.(Simple)
	if !ok{
		return false
	}

	return toCompare.Kind == t.Kind

}

func (array Array)Compare(datatype interface{})bool{
	toCompare, ok := datatype.(Array)
	if !ok{
		return false
	}

	if toCompare.Length != array.Length{
			return false
		}


	if reflect.TypeOf(array.Of) != reflect.TypeOf(toCompare.Of){
		return false
	}
	return Compare(array.Of, toCompare.Of)


}


func (pointer Pointer) Compare(datatype interface{}) bool{
	toCompare, ok := datatype.(Pointer)
	if !ok{
		return false
	}
	return Compare(pointer.PointsTo, toCompare.PointsTo)
}

func Compare(dataType1 interface{}, dataType2 interface{}) bool {
	switch dataType1.(type) {
	case Pointer:
		return dataType1.(Pointer).Compare(dataType2)
	case Array:
		return dataType1.(Array).Compare(dataType2)
	case Simple:
		return dataType1.(Simple).Compare(dataType2)

	default:
		panic(errorhandler.UnexpectedCompilerError())
	}
}

func Fmt(datatype interface{}) string{
	switch datatype.(type) {
	case Pointer:
		return "*"+ Fmt(datatype.(Pointer).PointsTo)
	case Array:
		array := datatype.(Array)
		return "["+strconv.Itoa(array.Length)+"]"+Fmt(array.Of)
	case Simple:
		simpleDataType :=  datatype.(Simple)
		if simpleDataType.Kind == KindByte{
			return "byte"
		}else{
			if simpleDataType.Kind == KindBool{
				return "bool"
			}else{
				return "void"
			}
		}

	default:
		panic(errorhandler.UnexpectedCompilerError())
	}
}
func NewFunction(returnDataType interface{}, argsDataType[] interface{}) Function {
	return Function{Return: returnDataType, Args: argsDataType}
}

func NewPointer(pointsTo interface{})Pointer{
	return Pointer{Size: 2, PointsTo: pointsTo}
}

func NewArray(length int, datatype interface{}) Array{
	return Array{Length: length, Of: datatype}
}

func NewBool() Simple {
	return Simple{Size: 1, Kind: KindBool}
}

func NewByte() Simple {
	return Simple{Size: 1, Kind: KindByte}
}

func NewVoid() Simple {
	return Simple{Size: 0, Kind: KindVoid}
}

func newSymbol(identifier string, datatype interface{})*Symbol{
	symbol := new(Symbol)
	symbol.Identifier = identifier
	switch datatype.(type){
	case Function:
		symbol.IsFunction = true
	default:
		symbol.IsFunction = false
	}
	symbol.DataType = datatype
	return symbol
}

func CreateGlobalScope()*Scope{
	return &Scope{
		SubScopes:        make([]*Scope, 0),
		NumberOfSubScope: 0,
		Parent:           nil,
		Symbols:          make(SymbolTable),
	}
}

func newScope(parent *Scope, parentSymbols SymbolTable) *Scope{
	return &Scope{
		SubScopes:        nil,
		NumberOfSubScope: 0,
		Parent:           parent,
		Symbols:          parentSymbols,
	}
}

func (scope *Scope)AddSubScope(){
	child := newScope(scope, scope.Symbols)
	scope.SubScopes = append(scope.SubScopes, child)
	scope.NumberOfSubScope += 1
}

func (scope *Scope )AddSymbol(identifier string, datatype interface{}) bool {
	_, exists := scope.Symbols[identifier]
	if exists{
		return false
	}
	symbol := newSymbol(identifier, datatype)
	scope.Symbols[identifier] = symbol
	return true
}

func GetSize(datatype interface{}) int{
	switch datatype.(type){
	case Pointer:
		return datatype.(Pointer).Size
	case Array:
		return datatype.(Array).SizeOfElements() * datatype.(Array).Length
	case Simple:
		return datatype.(Simple).Size
	default:
		panic(errorhandler.UnexpectedCompilerError())
	}

}

func IsAnArray(datatype interface{})bool{
	switch datatype.(type){
	case Array:
		return true

	default:
		return false
	}

}

func IsNumeric(datatype interface{})bool{
	switch datatype.(type){
	case Pointer:
		return true
	case Simple:
		return IsByte(datatype)
	default:
		return false
	}

}
func IsByte(datatype interface{})bool{
	switch datatype.(type){
	case Simple:
		if datatype.(Simple).Kind == KindByte{
			return true
		}else{
			return false
		}
	default:
		return false
	}

}