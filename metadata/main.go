package main

import (
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

const (
	FooType DataType = iota + 1
	BarType
)

// Shows how to use data type metadata.
func main() {
	bs := makeSlice()
	// If we don't know a concrete type of the encoded data, metadata will help
	// us.
	dt, _, err := UnmarshalDataType(bs)
	assert.EqualError(err, nil)

	switch dt {
	case FooType:
		// Unmarhsal Foo
	case BarType:
		// Unmarshal Bar
	default:
		panic("unexpected DataType")
	}
}

func makeSlice() (bs []byte) {
	var (
		foo  = Foo{a: 1, b: true}
		size = SizeMetaFoo(foo)
	)
	bs = make([]byte, size)
	MarshalMetaFoo(foo, bs)
	return
}
