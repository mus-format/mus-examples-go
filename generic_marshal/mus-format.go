package main

import (
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size functions.
// -----------------------------------------------------------------------------

func MarshalFooMUS(foo Foo, bs []byte) (n int) {
	return varint.MarshalInt(foo.num, bs)
}
func UnmarshalFooMUS(bs []byte) (foo Foo, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs[n:])
	return
}
func SizeFooMUS(foo Foo) (size int) {
	return varint.SizeInt(foo.num)
}
func SkipFooMUS(bs []byte) (n int, err error) {
	return varint.SkipInt(bs)
}

func MarshalBarMUS(bar Bar, bs []byte) (n int) {
	return ord.MarshalString(bar.str, nil, bs)
}
func UnmarshalBarMUS(bs []byte) (bar Bar, n int, err error) {
	bar.str, n, err = ord.UnmarshalString(nil, bs[n:])
	return
}
func SizeBarMUS(bar Bar) (size int) {
	return ord.SizeString(bar.str, nil)
}
func SkipBarMUS(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}
