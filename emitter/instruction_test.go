package emitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestI6XKK(t *testing.T) {
	type cases struct{
		description     string
		x 				byte
		kk 				byte
		expectedOpcode	Opcode
	}

	var opcode Opcode
	opcode[0]=0x62
	opcode[1]=0xAA
	testCases :=[]cases{
		{
			description: "0x62AA",
			x: 0x02,
			kk: 0xAA,
			expectedOpcode: opcode,
		},
	}
	for _, scenario := range testCases{

		t.Run(scenario.description, func(t *testing.T) {
			opcode := I6XKK(scenario.x, scenario.kk)
			assert.Equal(t, scenario.expectedOpcode, opcode)
		})
	}
}
func TestI1NNN(t *testing.T) {
	type cases struct{
		description     string
		nnn 			uint16
		expectedOpcode	Opcode
	}

	var opcode Opcode
	opcode[0]=0x12
	opcode[1]=0xAA
	testCases :=[]cases{
		{
			description: "0x12AA",
			nnn: 0x2AA,
			expectedOpcode: opcode,
		},
	}
	for _, scenario := range testCases{

		t.Run(scenario.description, func(t *testing.T) {
			opcode := I1NNN(scenario.nnn)
			assert.Equal(t, scenario.expectedOpcode, opcode)
		})
	}
}