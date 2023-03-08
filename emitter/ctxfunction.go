package emitter

type FunctionCtx struct{
	Registers       *RegistersGuide
	StackReferences *StackReferences
}

func NewCtxFunction(registers *RegistersGuide, stackReferences *StackReferences)*FunctionCtx {
	return &FunctionCtx{registers,stackReferences}
}
