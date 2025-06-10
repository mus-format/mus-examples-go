package main

import (
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

func main() {
	var (
		bs  []byte
		in  Instruction // Interface.
		err error
	)

	// Marshal Copy instruction.
	copy := Copy{start: 10, end: 20}
	bs = make([]byte, InstructionMUS.Size(copy))
	InstructionMUS.Marshal(copy, bs)

	// Unmarshal Copy instruction.
	in, _, err = InstructionMUS.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(in, copy)

	// Marshal Insert instruction.
	insert := Insert{str: "hello world"}
	bs = make([]byte, InstructionMUS.Size(insert))
	InstructionMUS.Marshal(insert, bs)

	// Unmarshal Insert instruction.
	in, _, err = InstructionMUS.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(in, insert)
}
