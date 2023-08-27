package main

import (
	"strconv"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	dvs "github.com/mus-format/mus-dvs-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Data Type Metadata (DTM).
// -----------------------------------------------------------------------------

const (
	FooV1DTM com.DTM = iota
	FooV2DTM
	BarV1DTM
)

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size functions.
// -----------------------------------------------------------------------------

// FooV1
func MarshalFooV1MUS(foo FooV1, bs []byte) (n int) {
	return varint.MarshalInt(foo.num, bs)
}

func UnmarshalFooV1MUS(bs []byte) (foo FooV1, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs)
	return
}

func SizeFooV1MUS(foo FooV1) (size int) {
	return varint.SizeInt(foo.num)
}

// FooV2
func MarshalFooV2MUS(foo FooV2, bs []byte) (n int) {
	return ord.MarshalString(foo.str, bs)
}

func UnmarshalFooV2MUS(bs []byte) (foo FooV2, n int, err error) {
	foo.str, n, err = ord.UnmarshalString(bs)
	return
}

func SizeFooV2MUS(foo FooV2) (size int) {
	return ord.SizeString(foo.str)
}

// BarV1
func MarshalBarV1MUS(foo BarV1, bs []byte) (n int) {
	return ord.MarshalBool(foo.b, bs)
}

func UnmarshalBarV1MUS(bs []byte) (foo BarV1, n int, err error) {
	foo.b, n, err = ord.UnmarshalBool(bs)
	return
}

func SizeBarV1MUS(foo BarV1) (size int) {
	return ord.SizeBool(foo.b)
}

// -----------------------------------------------------------------------------
// Data Type Metadata Support (DTS).
// -----------------------------------------------------------------------------

var FooV1DTS = dts.New[FooV1](FooV1DTM,
	mus.MarshallerFn[FooV1](MarshalFooV1MUS),
	mus.UnmarshallerFn[FooV1](UnmarshalFooV1MUS),
	mus.SizerFn[FooV1](SizeFooV1MUS),
)

var FooV2DTS = dts.New[FooV2](FooV2DTM,
	mus.MarshallerFn[FooV2](MarshalFooV2MUS),
	mus.UnmarshallerFn[FooV2](UnmarshalFooV2MUS),
	mus.SizerFn[FooV2](SizeFooV2MUS),
)

var BarV1DTS = dts.New[BarV1](BarV1DTM,
	mus.MarshallerFn[BarV1](MarshalBarV1MUS),
	mus.UnmarshallerFn[BarV1](UnmarshalBarV1MUS),
	mus.SizerFn[BarV1](SizeBarV1MUS),
)

// -----------------------------------------------------------------------------
// Data Versioning Support (DVS).
// -----------------------------------------------------------------------------

// First, we need to create a registry of all the versions we support. To do
// this, we use the dvs.Version type. This type is actually quite simple, it
// contains DTS and migrate functions, one of which migrates the old version to
// the current one, and the other do the opposite - migrates current to old.
//
// PLEASE NOTE that the index of each version in the registry must be equal
// to its DTM, i.e:
//
//	registry[FooV1DTM] == dvs.Version[FooV1, Foo]
//	registry[FooV2DTM] == dvs.Version[FooV2, Foo]
//	registry[BarV1DTM] == dvs.Version[BarV1, Foo]
//
// Thanks to this, we can very quickly get the version we need from the
// registry.
var registry = dvs.NewRegistry(
	[]dvs.TypeVersion{
		// FooV1 version.
		dvs.Version[FooV1, Foo]{
			DTS: FooV1DTS,
			MigrateOld: func(t FooV1) (v Foo, err error) {
				v.str = strconv.Itoa(t.num)
				return
			},
			MigrateCurrent: func(v Foo) (t FooV1, err error) {
				t.num, err = strconv.Atoi(v.str)
				return
			},
		},
		// FooV2 version.
		dvs.Version[FooV2, Foo]{
			DTS: FooV2DTS,
			MigrateOld: func(t FooV2) (v Foo, err error) {
				return Foo(t), nil
			},
			MigrateCurrent: func(v Foo) (t FooV2, err error) {
				return FooV2(v), nil
			},
		},
		// BarV1 version.
		dvs.Version[BarV1, Bar]{
			DTS: BarV1DTS,
			MigrateOld: func(t BarV1) (v Bar, err error) {
				return Bar(t), nil
			},
			MigrateCurrent: func(v Bar) (t BarV1, err error) {
				return BarV1(v), nil
			},
		},
	},
)

// And finally we create versioning support for our types. Please note that we
// use a single registry for all the DVS types.
var FooDVS = dvs.New[Foo](registry)
var BarDVS = dvs.New[Bar](registry)
