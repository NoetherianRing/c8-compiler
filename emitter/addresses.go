package emitter


type Stack struct {
	References    map[string]*Reference
	SubReferences []*Stack
}

type Reference struct{
	identifier      string
	positionInStack int
}

func NewStackReferences() *Stack {
	stackReferences := new(Stack)
	stackReferences.SubReferences = make([]*Stack,0)
	stackReferences.References = make(map[string]*Reference)
	return stackReferences
}

func (references *Stack) AddSubReferences() {
	subReference := NewStackReferences()
	for key, val := range references.References {
		subReference.References[key] = val
	}
	references.SubReferences = append(references.SubReferences, subReference)
}

func (references *Stack) AddReference(ident string, positionStack int) {
	references.References[ident] = &Reference{ident, positionStack}
}


func (references *Stack) GetReference(ident string) (*Reference, bool) {
	val, ok:= references.References[ident]
	return val, ok
}

