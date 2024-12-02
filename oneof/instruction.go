package main

import "fmt"

// Interface to Marshal/Unmarshal.
type Instruction interface {
	Do()
}

// Copy implements Instruction and MarshallerMUS interfaces.
type Copy struct {
	start int
	end   int
}

func (c Copy) Do() {
	fmt.Printf("copy from %v to %v\n", c.start, c.end)
}

func (c Copy) MarshalMUS(bs []byte) (n int) {
	return CopyDTS.Marshal(c, bs)
}

func (c Copy) SizeMUS() (size int) {
	return CopyDTS.Size(c)
}

// Insert implements Instruction and MarshallerMUS interfaces.
type Insert struct {
	str string
}

func (i Insert) Do() {
	fmt.Printf("insert '%v'\n", i.str)
}

func (i Insert) MarshalMUS(bs []byte) (n int) {
	return InsertDTS.Marshal(i, bs)
}

func (i Insert) SizeMUS() (size int) {
	return InsertDTS.Size(i)
}
