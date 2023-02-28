package emitter
const(
	MEMORY                   = 0xFFF - 0x1FF //ROM files starts at position 0x1FF
	AddressGlobalSection     = 6             //The memory position in which the section of global globalVariables starts
	StackSectionStartPointer = 10            //The memory position in which we store the address in which starts the stack section
	RegisterStackAddres1     = 4 			//The index of the register that saves the first 8 bits of the stack address
	RegisterStackAddres2     = 5			//The index of the register that saves the last 8 bits of the stack address

)
