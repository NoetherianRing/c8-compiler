package emitter


type Addresses struct {
	References   map[string]*Reference //by address we meant the position of the variable in the stack
	SubAddresses []*Addresses
}

type Reference struct{
	identifier  string
	positionStack int
}

func NewAddresses() *Addresses{
	addresses := new(Addresses)
	addresses.SubAddresses = make([]*Addresses,0)
	addresses.References = make(map[string]*Reference)
	return addresses
}

func (addresses Addresses) AddSubAddresses() {
	subAddress := NewAddresses()
	for key, val := range addresses.References {
		subAddress.References[key] = val
	}
	addresses.SubAddresses = append(addresses.SubAddresses, subAddress)
}

func (addresses Addresses) AddAddress(ident string, positionStack int) {
	addresses.References[ident] = &Reference{ident, positionStack}
}


func (addresses Addresses) GetReference(ident string) (*Reference, bool) {
	val, ok:= addresses.References[ident]
	return val, ok
}

