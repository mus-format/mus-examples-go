package main

import (
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

type Foo struct {
	a int
	b bool
}

// -----------------------------------------------------------------------------
func MarshalMetaFoo(foo Foo, bs []byte) (n int) {
	n = MarshalDataType(FooType, bs)
	return n + MarshalFoo(foo, bs[n:])
}

func SizeMetaFoo(foo Foo) (size int) {
	size += SizeDataType(FooType)
	return size + SizeFoo(foo)
}

// -----------------------------------------------------------------------------
func MarshalFoo(foo Foo, bs []byte) (n int) {
	n = varint.MarshalInt(foo.a, bs)
	return n + ord.MarshalBool(foo.b, bs[n:])
}

func UnmarshalFoo(bs []byte) (foo Foo, n int, err error) {
	foo.a, n, err = varint.UnmarshalInt(bs[n:])
	if err != nil {
		return
	}
	var n1 int
	foo.b, n1, err = ord.UnmarshalBool(bs[n:])
	n += n1
	return
}

func SizeFoo(foo Foo) (size int) {
	size += varint.SizeInt(foo.a)
	return size + ord.SizeBool(foo.b)
}
