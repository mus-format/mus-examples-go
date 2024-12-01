package main

import (
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

// Shows how to Marshal/Unmarshal interface.
func main() {
	var (
		bs    []byte
		instr Instruction // interface
		err   error
	)

	// Marshals Copy Instruction.
	copy := Copy{start: 10, end: 20}
	bs = make([]byte, SizeInstructionMUS(copy))
	MarshalInstructionMUS(copy, bs)

	// Unmarshals Copy Instruction.
	instr, _, err = UnmarshalInstructionMUS(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(instr, copy)

	// Marshals Insert Instruction.
	insert := Insert{str: "hello world"}
	bs = make([]byte, SizeInstructionMUS(insert))
	MarshalInstructionMUS(insert, bs)

	// Unmarshals Insert Instruction.
	instr, _, err = UnmarshalInstructionMUS(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(instr, insert)
}
