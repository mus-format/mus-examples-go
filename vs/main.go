package main

import (
	dtms "github.com/mus-format/mus-dtms-go"
	"github.com/ymz-ncnk/assert"
)

// This example demonstrates how to use the mus-vs-go module.
//
// mus-vs-go provides versioning support for the mus-go serializer. With the
// mus-vs-go module we can do 2 thigs:
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
	bs, _, _ = FooVS.MakeBSAndMarshalMUS(FooV1DTM, foo)

	// We should find the FooV1 version in bs.
	fooV1, _, _ := FooV1DTMS.UnmarshalMUS(bs)
	assert.EqualDeep(fooV1, FooV1{num: 5})

	// 2. Unmarshal the old version as if it was the current version.
	// Fill bs with the FooV1 version.
	fooV1 = FooV1{num: 10}
	bs = make([]byte, FooV1DTMS.SizeMUS(fooV1))
	FooV1DTMS.MarshalMUS(fooV1, bs)

	// Unmarshals the FooV1 version from bs, and then migrates it to the Foo type.
	dt, foo, _, _ := FooVS.UnmarshalMUS(bs)

	// We have to get the correct Foo.
	assert.Equal[dtms.DTM](dt, FooV1DTM)
	assert.EqualDeep(foo, Foo{str: "10"})
}
