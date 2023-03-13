package emitter
const(
	Memory	                   = 4096
	RomStart	               = 0x200		  //ROM files starts at position 0x200
	AddressGlobalSection       = RomStart + 8             //The memory position in which the section of global globalVariables starts
	RegisterStackAddress1 	   = 0xD            //The index of the register that saves the first 8 bits of the stack address
	RegisterStackAddress2 	   = 0xE            //The index of the register that saves the last 8 bits of the stack address
	AmountOfRegistersToOperate = 11            //The amount of registers that are allowed to use in a operation
	Carry					   = 0xF
	True                	   = 1
	False                 	   = 0
	SizePointer 			   = 2

	)
