package main

import (
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

// Shows how to implement the oneof feature. There is an Instruction interface,
// its Copy and Insert implementations.
func main() {
	var (
		bs    []byte
		instr Instruction // Interface.
		err   error
	)

	// Marshal Copy instruction.
	copy := Copy{start: 10, end: 20}
	bs = make([]byte, SizeInstructionMUS(copy))
	MarshalInstructionMUS(copy, bs)

	// Unmarshal Copy instruction.
	instr, _, err = UnmarshalInstructionMUS(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(instr, copy)

	// Marshal Insert instruction.
	insert := Insert{str: "hello world"}
	bs = make([]byte, SizeInstructionMUS(insert))
	MarshalInstructionMUS(insert, bs)

	// Unmarshal Insert instruction.
	instr, _, err = UnmarshalInstructionMUS(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(instr, insert)
}
