package emitter

type FunctionCtx struct{
	Stack *Stack
}

func NewCtxFunction(stackReferences *Stack)*FunctionCtx {
	return &FunctionCtx{stackReferences}
}
