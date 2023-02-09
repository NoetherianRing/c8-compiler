package emitter
const(
	MEMORY        = 0xFFF - 0x1FF        //ROM files starts at position 0x1FF
	GlobalSectionStart = 6 				//The memory position in which the section of global variables starts
	FirstRegisterForStackSection = 2	//we use the registers v2 and v3 to save the address of the stack section
	SecondRegisterForStackSection = 3


)
