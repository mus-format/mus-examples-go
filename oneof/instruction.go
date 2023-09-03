package main

import "fmt"

// Interface to Marshal/Unmarshal.
type Instruction interface {
	Do()
}

// Copy implements the Instruction interface.
type Copy struct {
	start int
	end   int
}

func (c Copy) Do() {
	fmt.Printf("copy from %v to %v\n", c.start, c.end)
}

// Insert implements the Instruction interface.
type Insert struct {
	str string
}

func (i Insert) Do() {
	fmt.Printf("insert '%v'\n", i.str)
}
