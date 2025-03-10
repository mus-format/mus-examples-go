package main

import (
	"fmt"
	"strconv"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// DTMs (Data Type Metadata.
const (
	FooV1DTM com.DTM = iota + 1
	FooV2DTM
)

// Serializers.
var (
	FooMUS   = fooMUS{}
	FooV1MUS = fooV1MUS{}
	FooV2MUS = fooV2MUS{}
)

// DTSs (Data Type metadata Support).
var (
	FooV1DTS = dts.New[FooV1](FooV1DTM, FooV1MUS)
	FooV2DTS = dts.New[FooV2](FooV2DTM, FooV2MUS)
)

// fooV1MUS implements mus.Serializer for FooV1.
type fooV1MUS struct{}

func (s fooV1MUS) Marshal(foo FooV1, bs []byte) (n int) {
	return varint.Int.Marshal(foo.num, bs)
}

func (s fooV1MUS) Unmarshal(bs []byte) (foo FooV1,
	n int, err error) {
	foo.num, n, err = varint.Int.Unmarshal(bs[n:])
	return
}

func (s fooV1MUS) Size(foo FooV1) (size int) {
	return varint.Int.Size(foo.num)
}

func (s fooV1MUS) Skip(bs []byte) (n int, err error) {
	return varint.Int.Skip(bs)
}

// fooV2MUS implements mus.Serializer for FooV2.
type fooV2MUS struct{}

func (s fooV2MUS) Marshal(foo FooV2, bs []byte) (n int) {
	return ord.String.Marshal(foo.str, bs)
}

func (s fooV2MUS) Unmarshal(bs []byte) (foo FooV2,
	n int, err error) {
	foo.str, n, err = ord.String.Unmarshal(bs[n:])
	return
}

func (s fooV2MUS) Size(foo FooV2) (size int) {
	return ord.String.Size(foo.str)
}

func (s fooV2MUS) Skip(bs []byte) (n int, err error) {
	return ord.String.Skip(bs)
}

// fooMUS implements mus.Serializer for Foo.
type fooMUS struct{}

func (s fooMUS) Marshal(foo Foo, bs []byte) (n int) {
	return FooV2DTS.Marshal(FooV2(foo), bs)
}

func (s fooMUS) Unmarshal(bs []byte) (foo Foo, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(bs)
	if err != nil {
		return
	}
	switch dtm {
	case FooV1DTM:
		var fooV1 FooV1
		fooV1, n, err = FooV1DTS.UnmarshalData(bs[n:])
		if err != nil {
			return
		}
		foo = migrateFooV1(fooV1)
	case FooV2DTM:
		var fooV2 FooV2
		fooV2, n, err = FooV2DTS.UnmarshalData(bs[n:])
		foo = Foo(fooV2)
	default:
		panic(fmt.Sprintf("unexpected %v DTM", dtm))
	}
	return
}

func (s fooMUS) Size(foo Foo) (size int) {
	return FooV2DTS.Size(FooV2(foo))
}

func (s fooMUS) Skip(bs []byte) (n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(bs)
	if err != nil {
		return
	}
	switch dtm {
	case FooV1DTM:
		return FooV1DTS.Skip(bs[n:])
	case FooV2DTM:
		return FooV2DTS.Skip(bs[n:])
	default:
		panic(fmt.Sprintf("unexpected %v DTM", dtm))
	}
}

func migrateFooV1(f FooV1) Foo {
	return Foo{
		str: strconv.Itoa(f.num),
	}
}
