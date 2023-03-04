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
	var i8xy0 Opcode
	i8xy0[0] =0x80 | (x)
	i8xy0[1] = 0x00 | (y << 4)
	return i8xy0
}

//I8XY1 writes in an Opcode the chip 8 instruction 8XY1 which set Vx = Vx | Vy
func I8XY1(x byte, y byte) Opcode {
	var i8xy1 Opcode
	i8xy1[0] =0x80 | (x)
	i8xy1[1] = 0x00 | (y << 4)
	return i8xy1
}

//I8XY2 writes in an Opcode the chip 8 instruction 8XY2 which set Vx = Vx & Vy
func I8XY2(x byte, y byte) Opcode {
	var i8xy2 Opcode
	i8xy2[0] =0x80 | (x)
	i8xy2[1] = 0x00 | (y << 4)
	return i8xy2
}

//I8XY3 writes in an Opcode the chip 8 instruction 8XY3 which set Vx = Vx ^ Vy
func I8XY3(x byte, y byte) Opcode {
	var i8xy3 Opcode
	i8xy3[0] =0x80 | (x)
	i8xy3[1] = 0x00 | (y << 4)
	return i8xy3
}
//I8XY4 writes in an Opcode the chip 8 instruction 8XY4 which set Vx = Vx + Vy
//If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func I8XY4(x byte, y byte) Opcode {
	var i8xy4 Opcode
	i8xy4[0] =0x80 | (x)
	i8xy4[1] = 0x00 | (y << 4)
	return i8xy4
}

//I8XY7 writes in an Opcode the chip 8 instruction 8XY7
//if Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func I8XY7(x byte, y byte)Opcode{
	var i8xy7 Opcode
	i8xy7[0] = 0x80 | x
	i8xy7[1] = 0x05 | (y << 4)
	return i8xy7
}

//I8XY5 writes in an Opcode the chip 8 instruction 8XY5 (Vx = Vx - Vy)
//if Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
func I8XY5(x byte, y byte)Opcode{
	var i8xy5 Opcode
	i8xy5[0] = 0x80 | x
	i8xy5[1] = 0x05 | (y << 4)
	return i8xy5
}
//I8XY6 writes in an Opcode the chip 8 instruction 8XY6 which set Vx = Vx >> 1
func I8XY6(x byte) Opcode {
	var i8xy6 Opcode
	i8xy6[0] =0x80 | (x)
	i8xy6[1] = 0x00
	return i8xy6
}


//I8XYE writes in an Opcode the chip 8 instruction 8XYE which set Vx = Vx << 1
func I8XYE(x byte) Opcode {
	var i8xyE Opcode
	i8xyE[0] =0x80 | (x)
	i8xyE[1] = 0x00
	return i8xyE
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
//I4XKK writes in an Opcode the chip 8 instruction 4XKK which skip the next instruction if vx != kk.
func I4XKK(x byte, kk byte)Opcode{
	var i4xkk Opcode
	i4xkk[0] = 0x40 | x
	i4xkk[1] = kk
	return i4xkk
}

//I2NNN writes in an Opcode the chip 8 instruction I2NNN CALL (ADDR)
func I2NNN(nnn uint16)Opcode{
	var i2nnn Opcode
	i2nnn[0] = 0x20 | byte(nnn >> 8)
	i2nnn[1] = byte(nnn)
	return i2nnn
}

//I00EE writes in an Opcode the chip 8 instruction 00EE which returns from a subroutine
func I00EE()Opcode{
	var i2nnn Opcode
	i2nnn[0] = 0x00
	i2nnn[1] = 0xEE
	return i2nnn
}

//I00E0 writes in an Opcode the chip 8 instruction 00E0 which cleans the monitor
func I00E0()Opcode{
	var i2nnn Opcode
	i2nnn[0] = 0x00
	i2nnn[1] = 0xE0
	return i2nnn
}

//I7XKK writes in an Opcode the chip 8 instruction 7XKK (vx = vx + kk)
func I7XKK(x byte, kk byte)Opcode{
	var i7xkk Opcode
	i7xkk[0] = 0x40 | x
	i7xkk[1] = kk
	return i7xkk
}

//IFX29 writes in an Opcode the chip 8 instruction FX29 which  set I = location of sprite for digit Vx.
func IFX29(x byte)Opcode{
	var ifx29 Opcode
	ifx29[0] = 0xf0 | x
	ifx29[1] = 0x29
	return ifx29
}


//IFX18 writes in an Opcode the chip 8 instruction FX18 which  set sound timer = Vx.
func IFX18(x byte)Opcode{
	var ifx18 Opcode
	ifx18[0] = 0xf0 | x
	ifx18[1] = 0x18
	return ifx18
}


//IFX15 writes in an Opcode the chip 8 instruction FX15 which set delay timer = Vx.
func IFX15(x byte)Opcode{
	var ifx15 Opcode
	ifx15[0] = 0xf0 | x
	ifx15[1] = 0x15
	return ifx15
}
