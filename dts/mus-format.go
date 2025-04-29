package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// DTM (Data Type Metadata) definitions.
const (
	FooDTM com.DTM = iota + 1
	BarDTM
)

// Serializers.
var (
	FooMUS = fooMUS{}
	BarMUS = barMUS{}
)

// DTS (Data Type metadata Support) definitions.
var (
	FooDTS = dts.New[Foo](FooDTM, FooMUS)
	BarDTS = dts.New[Bar](BarDTM, BarMUS)
)

// fooMUS implements mus.Serializer for Foo.
type fooMUS struct{}

func (s fooMUS) Marshal(foo Foo, bs []byte) (n int) {
	return varint.Int.Marshal(foo.num, bs)
}

func (s fooMUS) Unmarshal(bs []byte) (foo Foo, n int, err error) {
	foo.num, n, err = varint.Int.Unmarshal(bs[n:])
	return
}

func (s fooMUS) Size(foo Foo) (size int) {
	return varint.Int.Size(foo.num)
}

func (s fooMUS) Skip(bs []byte) (n int, err error) {
	return varint.Int.Skip(bs)
}

// barMUS implements mus.Serializer for Bar.
type barMUS struct{}

func (s barMUS) Marshal(bar Bar, bs []byte) (n int) {
	return ord.String.Marshal(bar.str, bs)
}

func (s barMUS) Unmarshal(bs []byte) (bar Bar, n int, err error) {
	bar.str, n, err = ord.String.Unmarshal(bs[n:])
	return
}

func (s barMUS) Size(bar Bar) (size int) {
	return ord.String.Size(bar.str)
}

func (s barMUS) Skip(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}
