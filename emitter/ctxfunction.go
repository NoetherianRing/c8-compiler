package emitter

type FunctionCtx struct{
	Registers *RegistersGuide
	Addresses *Addresses
}

func NewCtxFunction(registers *RegistersGuide, addresses *Addresses)*FunctionCtx {
	return &FunctionCtx{registers,addresses}
}
