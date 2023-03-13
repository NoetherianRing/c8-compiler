package emitter

const NonAvailable = -1

//ResultRegIndex stores the indexes of the registers in which a result is stored, and a field that says if the result
//stored in the registers is a pointer (or a simple if false)
type ResultRegIndex struct{
	highBitsIndex byte //if the result is a pointer, highBitsIndex is the index of the register that stores the first 8 bits.
	lowBitsIndex  byte //if the result is a pointer, lowBitsIndex is the index of the register  that stores the last 8 bits.
	//if the result is a simple, lowBitsIndex is the index of the register in which the entire result is stored.
	isPointer	bool
}

//RegisterHandler handle the availability of registers
type RegisterHandler struct{
	available [AmountOfRegistersToOperate]bool
	nextAvailableRegister int
}

func NewRegisterHandler() *RegisterHandler{
	registerHandler := new(RegisterHandler)
	for i:=0; i<AmountOfRegistersToOperate; i++{
		registerHandler.available[i] = true

	}
	registerHandler.nextAvailableRegister = 0
	return registerHandler
}

//AllocSimple allocates a register for a simple, and returns its index. Returns false only when there are not enough registers
//available
func (handler *RegisterHandler) AllocSimple()(*ResultRegIndex, bool){
	index, ok := handler.alloc()
	return &ResultRegIndex{lowBitsIndex: index, isPointer: false}, ok
}

//AllocPointer allocates a register for a pointer, and returns its index. Returns false only when there are not enough registers
//available
func (handler *RegisterHandler) AllocPointer()(*ResultRegIndex, bool){
	highBitsIndex, ok := handler.alloc()
	if !ok{
		return nil, ok
	}
	lowBitsIndex, ok := handler.alloc()
	return &ResultRegIndex{highBitsIndex: highBitsIndex, lowBitsIndex: lowBitsIndex, isPointer: true}, ok

}

func (handler *RegisterHandler) alloc()(byte, bool){
	if handler.nextAvailableRegister != NonAvailable{
		registerToReturn := handler.nextAvailableRegister
		handler.available[registerToReturn] = false
		for i:=0; i<AmountOfRegistersToOperate; i++{
			if handler.available[i]{
				handler.nextAvailableRegister = i
				return byte(registerToReturn+2), true //we return from index 2 because v0 and v1 are reserved
			}
		}
		handler.nextAvailableRegister = NonAvailable
		return byte(registerToReturn+2), true
	}
	return 0, false

}

func (handler *RegisterHandler) Free(resultRegIndex *ResultRegIndex){
	if resultRegIndex.isPointer{
		handler.free(resultRegIndex.highBitsIndex)
	}
	handler.free(resultRegIndex.lowBitsIndex)

}

func (handler *RegisterHandler) free(index byte){
	handler.available[index] = true
	if handler.nextAvailableRegister > int(index){
		handler.nextAvailableRegister = int(index)
	}
}

func (handler *RegisterHandler) reserveRegister(index byte) bool{
	if handler.nextAvailableRegister == int(index){
		_, ok := handler.alloc()
		return ok
	}else{
		if !handler.available[index]{
			return false
		}else{
			handler.available[index] = false
			return true
		}
	}
}