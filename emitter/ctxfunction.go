package emitter

type FunctionCtx struct{
	Registers *RegistersGuide
	Addresses *StackReferences
}

func NewCtxFunction(registers *RegistersGuide, addresses *StackReferences)*FunctionCtx {
	return &FunctionCtx{registers,addresses}
}
