package emitter
const(
	MEMORY                = 0xFFF - 0x1FF //ROM files starts at position 0x1FF
	AddressGlobalSection  = 6             //The memory position in which the section of global globalVariables starts
	RegisterStackAddress1 = 14            //The index of the register that saves the first 8 bits of the stack address
	RegisterStackAddress2 = 15            //The index of the register that saves the last 8 bits of the stack address
	True                  = 1
	False                 = 0
)
