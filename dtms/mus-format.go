package main

import (
	dtms "github.com/mus-format/mus-dtms-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Data Type Metadata (DTM).
// -----------------------------------------------------------------------------

const (
	FooDTM dtms.DTM = iota
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
// Data Type Metadata Support (DTMS).
// -----------------------------------------------------------------------------

var FooDTMS = dtms.New[Foo](FooDTM,
	mus.MarshallerFn[Foo](MarshalFooMUS),
	mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
	mus.SizerFn[Foo](SizeFooMUS),
)
