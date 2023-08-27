package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Data Type Metadata (DTM).
// -----------------------------------------------------------------------------

const (
	FooDTM com.DTM = iota
	BarDTM
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

// -----------------------------------------------------------------------------
// Data Type Metadata Support (DTS).
// -----------------------------------------------------------------------------

var FooDTS = dts.New[Foo](FooDTM,
	mus.MarshallerFn[Foo](MarshalFooMUS),
	mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
	mus.SizerFn[Foo](SizeFooMUS),
)
