package main

import (
	com "github.com/mus-format/common-go"
	"github.com/ymz-ncnk/assert"
)

// This example demonstrates how to use mus-dvs-go.
//
// mus-dvs-go provides data versioning support for the mus-go serializer. With
// mus-dvs-go we can do 2 thigs:
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
	bs, _, _ = FooDVS.MakeBSAndMarshal(FooV1DTM, foo)

	// We should find the FooV1 version in bs.
	fooV1, _, _ := FooV1DTS.Unmarshal(bs)
	assert.EqualDeep(fooV1, FooV1{num: 5})

	// 2. Unmarshal the old version as if it was the current version.
	// Fill bs with the FooV1 version.
	fooV1 = FooV1{num: 10}
	bs = make([]byte, FooV1DTS.Size(fooV1))
	FooV1DTS.Marshal(fooV1, bs)

	// Unmarshals the FooV1 version from bs, and then migrates it to the current
	// version.
	dt, foo, _, _ := FooDVS.Unmarshal(bs)

	// We have to get the correct Foo.
	assert.Equal[com.DTM](dt, FooV1DTM)
	assert.EqualDeep(foo, Foo{str: "10"})
}
