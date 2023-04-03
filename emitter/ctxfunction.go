package emitter

type FunctionCtx struct {
	registerHandler *RegisterHandler
	stack           *Stack
}

func NewCtxFunction(registerHandler *RegisterHandler, stackReferences *Stack) *FunctionCtx {
	return &FunctionCtx{registerHandler, stackReferences}
}
