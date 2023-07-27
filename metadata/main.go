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
	wantFoo, bs := makeFooSlice()
	wantBar, _ := makeBarSlice()
	// If we don't know a concrete type of the encoded data, metadata will help
	// us.
	dt, n, err := UnmarshalDataType(bs)
	assert.EqualError(err, nil)

	switch dt {
	case FooType:
		// Unmarhsal Foo
		foo, _, err := UnmarshalFoo(bs[n:])
		assert.EqualError(err, nil)
		assert.EqualDeep(foo, wantFoo)
	case BarType:
		// Unmarshal Bar
		bar, _, err := UnmarshalBar(bs[n:])
		assert.EqualError(err, nil)
		assert.EqualDeep(bar, wantBar)
	default:
		panic("unexpected DataType")
	}
}

func makeFooSlice() (foo Foo, bs []byte) {
	foo = Foo{a: 1, b: true}
	bs = make([]byte, SizeMetaFoo(foo))
	MarshalMetaFoo(foo, bs)
	return
}

func makeBarSlice() (bar Bar, bs []byte) {
	bar = Bar{a: "hello world", b: true}
	bs = make([]byte, SizeMetaBar(bar))
	MarshalMetaBar(bar, bs)
	return
}
