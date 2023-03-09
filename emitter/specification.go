package emitter
const(
	Memory	              =  4096
	RomStart	          =  0x200		  //ROM files starts at position 0x200
	AddressGlobalSection  = RomStart + 6             //The memory position in which the section of global globalVariables starts
	RegisterStackAddress1 = 14            //The index of the register that saves the first 8 bits of the stack address
	RegisterStackAddress2 = 15            //The index of the register that saves the last 8 bits of the stack address
	True                  = 1
	False                 = 0
)
