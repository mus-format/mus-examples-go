package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// DTM (Data Type Metadata).
// -----------------------------------------------------------------------------

const (
	FooV1DTM com.DTM = iota + 1
	FooV2DTM
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size/Skip functions.
// -----------------------------------------------------------------------------

// FooV1

func MarshalFooV1MUS(foo FooV1, bs []byte) (n int) {
	return varint.MarshalInt(foo.num, bs)
}
func UnmarshalFooV1MUS(bs []byte) (foo FooV1, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs[n:])
	return
}
func SizeFooV1MUS(foo FooV1) (size int) {
	return varint.SizeInt(foo.num)
}
func SkipFooV1MUS(bs []byte) (n int, err error) {
	return varint.SkipInt(bs)
}

// FooV2

func MarshalFooV2MUS(foo FooV2, bs []byte) (n int) {
	return ord.MarshalString(foo.str, nil, bs)
}
func UnmarshalFooV2MUS(bs []byte) (foo FooV2, n int, err error) {
	foo.str, n, err = ord.UnmarshalString(nil, bs[n:])
	return
}
func SizeFooV2MUS(foo FooV2) (size int) {
	return ord.SizeString(foo.str, nil)
}
func SkipFooV2MUS(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}

// -----------------------------------------------------------------------------
// DTS (Data Type Metadata Support).
// -----------------------------------------------------------------------------

var FooV1DTS = dts.New[FooV1](FooV1DTM,
	mus.MarshallerFn[FooV1](MarshalFooV1MUS),
	mus.UnmarshallerFn[FooV1](UnmarshalFooV1MUS),
	mus.SizerFn[FooV1](SizeFooV1MUS),
	mus.SkipperFn(SkipFooV1MUS),
)
var FooV2DTS = dts.New[FooV2](FooV2DTM,
	mus.MarshallerFn[FooV2](MarshalFooV2MUS),
	mus.UnmarshallerFn[FooV2](UnmarshalFooV2MUS),
	mus.SizerFn[FooV2](SizeFooV2MUS),
	mus.SkipperFn(SkipFooV2MUS),
)
