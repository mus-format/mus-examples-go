package main

import (
	dtms "github.com/mus-format/mus-dtms-go"
	"github.com/ymz-ncnk/assert"
)

// This example demonstrates how to use the mus-dtms-go module.
//
// The mus-dtms-go module allows us to create data type metadata support for a
// type. Thus, for example, for the Foo type, we can create FooDTMS, which
// in turn allows us to encode/decode FooDTM + Foo data.
func main() {
	var (
		foo = Foo{num: 10}
		bs  = make([]byte, FooDTMS.SizeMUS(foo))
	)
	FooDTMS.MarshalMUS(foo, bs)
	// After marshal we are expecting to receive the following bs, where
	// {0} - FooDTM, {20} - Foo data.
	assert.EqualBytes(bs, []byte{0, 20})

	// Unmarshal should return the same foo.
	afoo, _, _ := FooDTMS.UnmarshalMUS(bs)
	assert.EqualDeep(afoo, foo)

	// And if the encoded DTM in bs is not FooDTM, we will receive
	// dtms.ErrWrongDTM.
	bs[0] = byte(BarDTM)
	_, _, err := FooDTMS.UnmarshalMUS(bs)
	assert.EqualError(err, dtms.ErrWrongDTM)
}
