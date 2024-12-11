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
	FooDTM com.DTM = iota + 1
	BarDTM
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size/Skip functions.
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

// -----------------------------------------------------------------------------
// DTS (Data Type Metadata Support).
// -----------------------------------------------------------------------------

var FooDTS = dts.New[Foo](FooDTM,
	mus.MarshallerFn[Foo](MarshalFooMUS),
	mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
	mus.SizerFn[Foo](SizeFooMUS),
	mus.SkipperFn(SkipFooMUS),
)
var BarDTS = dts.New[Bar](BarDTM,
	mus.MarshallerFn[Bar](MarshalBarMUS),
	mus.UnmarshallerFn[Bar](UnmarshalBarMUS),
	mus.SizerFn[Bar](SizeBarMUS),
	mus.SkipperFn(SkipBarMUS),
)
