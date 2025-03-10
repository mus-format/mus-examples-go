package main

import (
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example demonstrates how to implement the "oneof" feature. It includes
// an Instruction interface with its Copy and Insert implementations.
func main() {
	var (
		bs  []byte
		in  Instruction // Interface.
		err error
	)

	// Marshal Copy instruction.
	copy := Copy{start: 10, end: 20}
	bs = make([]byte, InstructionSer.Size(copy))
	InstructionSer.Marshal(copy, bs)

	// Unmarshal Copy instruction.
	in, _, err = InstructionSer.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(in, copy)

	// Marshal Insert instruction.
	insert := Insert{str: "hello world"}
	bs = make([]byte, InstructionSer.Size(insert))
	InstructionSer.Marshal(insert, bs)

	// Unmarshal Insert instruction.
	in, _, err = InstructionSer.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(in, insert)
}
