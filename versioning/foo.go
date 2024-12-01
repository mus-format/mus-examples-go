package main

import (
	"fmt"
	"strconv"

	dts "github.com/mus-format/mus-dts-go"
)

// Indicates the current version.
type Foo FooV2

func MarshalFooMUS(foo Foo, bs []byte) (n int) {
	return FooV2DTS.Marshal(FooV2(foo), bs)
}

func UnmarshalFooMUS(bs []byte) (foo Foo, n int, err error) {
	dtm, n, err := dts.UnmarshalDTM(bs)
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

func SizeFooMUS(foo Foo) (size int) {
	return FooV2DTS.Size(FooV2(foo))
}

func SkipFooMUS(bs []byte) (n int, err error) {
	dtm, n, err := dts.UnmarshalDTM(bs)
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
