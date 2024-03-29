package emitter

import "github.com/NoetherianRing/c8-compiler/symboltable"

const NonAvailable = -1

//ResultRegIndex stores the indexes of the registers in which a result is stored, and a field that says if the result
//stored in the registers is a pointer (or a simple if false)
type ResultRegIndex struct {
	highBitsIndex byte //if the result is a pointer, highBitsIndex is the index of the register that stores the first 8 bits.
	lowBitsIndex  byte //if the result is a pointer, lowBitsIndex is the index of the register  that stores the last 8 bits.
	//if the result is a simple, lowBitsIndex is the index of the register in which the entire result is stored.
	isPointer bool
}

//RegisterHandler handle the availability of registers
type RegisterHandler struct {
	available             [AmountOfRegistersToOperate - 2]bool //we subtract 2 because the register v0 and v1 are not handled in the same way
	nextAvailableRegister int
}

func NewRegisterHandler() *RegisterHandler {
	registerHandler := new(RegisterHandler)
	for i := 0; i < AmountOfRegistersToOperate-2; i++ {
		registerHandler.available[i] = true

	}
	registerHandler.nextAvailableRegister = 0
	return registerHandler
}

//Alloc receive a datatype and allocates registers to save its value, and returns its index.
//Returns false when there are not enough registers available or if the datatype is not valid
func (handler *RegisterHandler) Alloc(datatype interface{}) (*ResultRegIndex, bool) {
	switch datatype.(type) {
	case symboltable.Simple:
		return handler.AllocSimple()
	case symboltable.Pointer:
		return handler.AllocPointer()
	default:
		return nil, false
	}

}

//AllocSimple allocates a register for a simple, and returns its index. Returns false only when there are not enough registers
//available
func (handler *RegisterHandler) AllocSimple() (*ResultRegIndex, bool) {
	index, ok := handler.alloc()
	return &ResultRegIndex{lowBitsIndex: index, isPointer: false}, ok
}

//AllocPointer allocates two register for a pointer, and returns its index. Returns false only when there are not enough registers
//available
func (handler *RegisterHandler) AllocPointer() (*ResultRegIndex, bool) {
	highBitsIndex, ok := handler.alloc()
	if !ok {
		return nil, ok
	}
	lowBitsIndex, ok := handler.alloc()
	return &ResultRegIndex{highBitsIndex: highBitsIndex, lowBitsIndex: lowBitsIndex, isPointer: true}, ok

}

//alloc allocates a register and returns its index. Returns false only when there are not enough registers
func (handler *RegisterHandler) alloc() (byte, bool) {
	if handler.nextAvailableRegister != NonAvailable {
		registerToReturn := handler.nextAvailableRegister
		handler.available[registerToReturn] = false
		for i := 0; i < AmountOfRegistersToOperate-2; i++ {
			if handler.available[i] {
				handler.nextAvailableRegister = i
				return byte(registerToReturn + 2), true //we return from index 2 because v0 and v1 are reserved
			}
		}
		handler.nextAvailableRegister = NonAvailable
		return byte(registerToReturn + 2), true
	}
	return 0, false

}


func (handler *RegisterHandler) Free(resultRegIndex *ResultRegIndex) {
	if resultRegIndex.isPointer {
		handler.free(resultRegIndex.highBitsIndex - 2)
	}
	handler.free(resultRegIndex.lowBitsIndex - 2)

}

func (handler *RegisterHandler) free(index byte) {
	handler.available[index] = true
	if handler.nextAvailableRegister > int(index) ||
		handler.nextAvailableRegister == NonAvailable {
		handler.nextAvailableRegister = int(index)
	}
}
