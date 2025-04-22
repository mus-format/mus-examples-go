package main

import (
	"fmt"

	ext "github.com/mus-format/ext-mus-go"
)

// This example demonstrates how to implement the generic marshal function.
func main() {
	// Both Foo and Bar types implement ext.MarshallerMUS interface.

	bs := ext.MarshalMUS(Foo{num: 10}) // ext.MarshalMUS can be used with Foo ...
	fmt.Println(bs)

	bs = ext.MarshalMUS(Bar{str: "10"}) // ... and with Bar.
	fmt.Println(bs)
}
