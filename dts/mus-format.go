package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
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
	FooSer = fooSer{}
	BarSer = barSer{}
)

// DTS (Data Type metadata Support) definitions.
var (
	FooDTS = dts.New[Foo](FooDTM, FooSer)
	BarDTS = dts.New[Bar](BarDTM, BarSer)
)

// fooSer implements mus.Serializer for Foo.
type fooSer struct{}

func (s fooSer) Marshal(foo Foo, bs []byte) (n int) {
	return varint.Int.Marshal(foo.num, bs)
}

func (s fooSer) Unmarshal(bs []byte) (foo Foo, n int, err error) {
	foo.num, n, err = varint.Int.Unmarshal(bs[n:])
	return
}

func (s fooSer) Size(foo Foo) (size int) {
	return varint.Int.Size(foo.num)
}

func (s fooSer) Skip(bs []byte) (n int, err error) {
	return varint.Int.Skip(bs)
}

// barSer implements mus.Serializer for Bar.
type barSer struct{}

func (s barSer) Marshal(bar Bar, bs []byte) (n int) {
	return ord.String.Marshal(bar.str, bs)
}

func (s barSer) Unmarshal(bs []byte) (bar Bar, n int, err error) {
	bar.str, n, err = ord.String.Unmarshal(bs[n:])
	return
}

func (s barSer) Size(bar Bar) (size int) {
	return ord.String.Size(bar.str)
}

func (s barSer) Skip(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}
