package main

import "fmt"

// This example demonstrates how to implement the generic marshal function.
func main() {
	// Both Foo and Bar types implement MarshallerMUS interface.

	bs := MarshalMUS(Foo{num: 10}) // Can be used with Foo ...
	fmt.Println(bs)

	bs = MarshalMUS(Bar{str: "10"}) // ... and with Bar.
	fmt.Println(bs)
}

// Generic marshal function.
func MarshalMUS(v MarshallerMUS) (bs []byte) {
	bs = make([]byte, v.SizeMUS())
	v.MarshalMUS(bs)
	return
}
