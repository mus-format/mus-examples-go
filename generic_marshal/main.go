package main

import "fmt"

type MarshallerMUS interface {
	MarshalMUS(bs []byte) (n int)
	SizeMUS() (size int)
}

// Generic marshal function.
func MarshalMUS[T MarshallerMUS](t T) (bs []byte) {
	bs = make([]byte, t.SizeMUS())
	t.MarshalMUS(bs)
	return
}

// Demonstrates how to implement generic marshal function.
func main() {
	// Both Foo and Bar types implement MarshallerMUS interface.

	foo := Foo{num: 10}
	fn(foo) // Can be used with Foo ...

	bar := Bar{str: "10"}
	fn(bar) // ... and with Bar.
}

func fn[T MarshallerMUS](t T) {
	bs := MarshalMUS(t)
	fmt.Println(bs)
}
