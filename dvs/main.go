package main

import (
	com "github.com/mus-format/common-go"
	"github.com/ymz-ncnk/assert"
)

// This example demonstrates how to use the mus-dvs-go module.
//
// mus-dvs-go provides data versioning support for the mus-go serializer. With
// the mus-dvs-go module we can do 2 thigs:
// 1. Marshal the current version as if it was an old version.
// 2. Unmarshal the old version as if it was the current version.
func main() {
	var (
		foo Foo
		bs  []byte
	)
	// 1. Marshal the current version as if it was an old version.
	foo = Foo{str: "5"}
	// Migrates foo to the FooV1 version, and then marshals it to bs.
	bs, _, _ = FooDVS.MakeBSAndMarshalMUS(FooV1DTM, foo)

	// We should find the FooV1 version in bs.
	fooV1, _, _ := FooV1DTS.UnmarshalMUS(bs)
	assert.EqualDeep(fooV1, FooV1{num: 5})

	// 2. Unmarshal the old version as if it was the current version.
	// Fill bs with the FooV1 version.
	fooV1 = FooV1{num: 10}
	bs = make([]byte, FooV1DTS.SizeMUS(fooV1))
	FooV1DTS.MarshalMUS(fooV1, bs)

	// Unmarshals the FooV1 version from bs, and then migrates it to the current
	// version.
	dt, foo, _, _ := FooDVS.UnmarshalMUS(bs)

	// We have to get the correct Foo.
	assert.Equal[com.DTM](dt, FooV1DTM)
	assert.EqualDeep(foo, Foo{str: "10"})
}
