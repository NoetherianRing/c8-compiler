package emitter

type Opcode [2]byte

func I6XKK(x byte, kk byte) Opcode {
	var i6xkk Opcode
	i6xkk[0] =0x60 | x
	i6xkk[1] = kk
	return i6xkk
}
func I1NNN(nnn uint16) Opcode {
	var i1nnn Opcode
	i1nnn[0] =0x10 | byte(nnn >> 8)
	i1nnn[1] = byte(nnn)
	return i1nnn
}
