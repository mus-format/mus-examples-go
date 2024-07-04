package main

import (
	dts "github.com/mus-format/mus-dts-go"
	"github.com/ymz-ncnk/assert"
)

// This example demonstrates how to use mus-dts-go.
//
// mus-dts-go allows us to create data type metadata support for a type. Thus,
// for example, for the Foo type, we can create FooDTS, which in turn allows us
// to encode/decode FooDTM + Foo data.
func main() {
	var (
		foo = Foo{num: 10}
		bs  = make([]byte, FooDTS.Size(foo))
	)
	FooDTS.Marshal(foo, bs)
	// After marshal we are expecting to receive the following bs, where
	// {0} - FooDTM, {20} - Foo data.
	assert.EqualBytes(bs, []byte{0, 20})

	// Unmarshal should return the same foo.
	afoo, _, _ := FooDTS.Unmarshal(bs)
	assert.EqualDeep(afoo, foo)

	// And if the encoded DTM in bs is not FooDTM, we will receive
	// dts.ErrWrongDTM.
	bs[0] = byte(BarDTM)
	_, _, err := FooDTS.Unmarshal(bs)
	assert.EqualError(err, dts.ErrWrongDTM)
}
