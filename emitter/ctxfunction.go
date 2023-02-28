package emitter

type FunctionCtx struct{
	Registers *Registers
	Addresses *Addresses
}

func NewCtxFunction(registers *Registers, addresses *Addresses)*FunctionCtx {
	return &FunctionCtx{registers,addresses}
}
