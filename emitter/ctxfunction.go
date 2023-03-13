package emitter

type FunctionCtx struct{
	registerHandler *RegisterHandler
	stack           *Stack
}

func NewCtxFunction(stackReferences *Stack)*FunctionCtx {
	registerHandler := NewRegisterHandler()
	return &FunctionCtx{registerHandler, stackReferences}
}
