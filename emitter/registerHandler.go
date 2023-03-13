package emitter

const NonAvailable = -1
type registersIndex [2]byte

func (index registersIndex) GetSize()int{
	return len(index)
}
func (index registersIndex) GetFirstHalfIndex()byte{
	return index[0]
}
func (index registersIndex) GetSecondHalfIndex()byte{
	return index[1]
}
func (index registersIndex) GetIndex()byte{
	return index.GetFirstHalfIndex()
}
func (index registersIndex) SetFirstHalfIndex(i byte){
	index[0] = i
}
func (index registersIndex) SetSecondHalfIndex(i byte){
	index[1] = i
}
func (index registersIndex) SetIndex(i byte){
	 index.SetFirstHalfIndex(i)
}



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

func (handler *RegisterHandler) TakeRegister() (byte, bool){
	if handler.nextAvailableRegister != NonAvailable{
		registerToReturn := handler.nextAvailableRegister
		handler.available[registerToReturn] = false
		for i:=0; i<AmountOfRegistersToOperate; i++{
			if handler.available[i]{
				handler.nextAvailableRegister = i
				return byte(registerToReturn), true
			}
		}
		handler.nextAvailableRegister = NonAvailable
		return byte(registerToReturn), true
	}
	return 0, false
}

func (handler *RegisterHandler) freeRegister(index byte){
	handler.available[index] = true
	if handler.nextAvailableRegister > int(index){
		handler.nextAvailableRegister = int(index)
	}
}