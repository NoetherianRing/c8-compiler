package emitter


type StackReferences struct {
	References   map[string]*Reference
	SubAddresses []*StackReferences
}

type Reference struct{
	identifier  string
	positionStack int
}

func NewStackReferences() *StackReferences {
	stackReferences := new(StackReferences)
	stackReferences.SubAddresses = make([]*StackReferences,0)
	stackReferences.References = make(map[string]*Reference)
	return stackReferences
}

func (addresses StackReferences) AddSubAddresses() {
	subAddress := NewStackReferences()
	for key, val := range addresses.References {
		subAddress.References[key] = val
	}
	addresses.SubAddresses = append(addresses.SubAddresses, subAddress)
}

func (addresses StackReferences) AddReference(ident string, positionStack int) {
	addresses.References[ident] = &Reference{ident, positionStack}
}


func (addresses StackReferences) GetReference(ident string) (*Reference, bool) {
	val, ok:= addresses.References[ident]
	return val, ok
}

