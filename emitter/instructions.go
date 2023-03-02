package emitter

type Opcode [2]byte

//I6XKK writes in an Opcode the chip 8 instruction 6XKK which set vx = kk
func I6XKK(x byte, kk byte) Opcode {
	var i6xkk Opcode
	i6xkk[0] =0x60 | x
	i6xkk[1] = kk
	return i6xkk
}
//I8XY0 writes in an Opcode the chip 8 instruction 8XY0 which set Vx = Vy
func I8XY0(x byte, y byte) Opcode {
	var i6xkk Opcode
	i6xkk[0] =0x80 | (x)
	i6xkk[1] = 0x00 | (y << 4)
	return i6xkk
}

//I1NNN writes in an Opcode the chip 8 instruction 1NNN which jumps to location nnn
func I1NNN(nnn uint16) Opcode {
	var i1nnn Opcode
	i1nnn[0] =0x10 | byte(nnn >> 8)
	i1nnn[1] = byte(nnn)
	return i1nnn
}
//IANNN writes in an Opcode the chip 8 instruction ANNN which set I = NNN
func IANNN(nnn uint16) Opcode {
	var iannn Opcode
	iannn[0] =0xA0 | byte(nnn >> 8)
	iannn[1] = byte(nnn)
	return iannn
}

//IFX55 writes in an Opcode the chip 8 instruction FX55 which stores registers V0 through Vx in memory starting at location I.
func IFX55(x byte) Opcode{
	var ifx55 Opcode
	ifx55[0] =0xF0 | x
	ifx55[1] = 0x55
	return ifx55
}
//IFX65 writes in an Opcode the chip 8 instruction FX65 which reads registers V0 through Vx from memory starting at location I.
func IFX65(x byte) Opcode{
	var ifx65 Opcode
	ifx65[0] =0xF0 | x
	ifx65[1] = 0x65
	return ifx65
}


//TODO: Borrar este mensaje, esto es NUXY
//IAXY0 writes in an Opcode the new chip 8 instruction AXY0 which was added in order of make pointers possible in this compiler
//It set I = (V1 << 8) | V2
func IAXY0(x byte, y byte)Opcode{
	var iaxy0 Opcode
	iaxy0[0] = 0xA0 | x
	iaxy0[1] = y
	return iaxy0
}

//IBXY0 writes in an Opcode the new chip 8 instruction BXY0 which was added in order of make pointers possible in this compiler
//It set Vx = byte(I>>8) Vy = byte(I)
func IBXY0(x byte, y byte)Opcode{
	var ibxy0 Opcode
	ibxy0[0] = 0xB0 | x
	ibxy0[1] = y
	return ibxy0
}

//IFX1E writes in an Opcode the chip 8 instruction FX1E which sets I = I + Vx.
func IFX1E(x byte)Opcode{
	var ifx1e Opcode
	ifx1e[0] = 0xF0 | x
	ifx1e[1] = 0x1E
	return ifx1e
}


//I3XKK writes in an Opcode the chip 8 instruction 3XKK which skip the next instruction if vx = kk.
func I3XKK(x byte, kk byte)Opcode{
	var i3xkk Opcode
	i3xkk[0] = 0x30 | x
	i3xkk[1] = kk
	return i3xkk
}