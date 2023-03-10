package emitter


type StackReferences struct {
	References    map[string]*Reference
	SubReferences []*StackReferences
}

type Reference struct{
	identifier  string
	positionStack int
}

func NewStackReferences() *StackReferences {
	stackReferences := new(StackReferences)
	stackReferences.SubReferences = make([]*StackReferences,0)
	stackReferences.References = make(map[string]*Reference)
	return stackReferences
}

func (references *StackReferences) AddSubReferences() {
	subReference := NewStackReferences()
	for key, val := range references.References {
		subReference.References[key] = val
	}
	references.SubReferences = append(references.SubReferences, subReference)
}

func (references *StackReferences) AddReference(ident string, positionStack int) {
	references.References[ident] = &Reference{ident, positionStack}
}


func (references *StackReferences) GetReference(ident string) (*Reference, bool) {
	val, ok:= references.References[ident]
	return val, ok
}

