package main

import (
	"fmt"
	"math/rand"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
)

// This example demonstrates how mus-dts-go can be used. Also, take a look at
// the generic_marshal example.
func main() {
	// Let's make a random data
	bs := randomData()
	// and Unmarshal DTM.
	dtm, n, err := dts.UnmarshalDTM(bs)
	if err != nil {
		panic(err)
	}
	// Now we can deserialize and process data depending on the DTM, which in
	// general allows us to receive data of different types.
	switch dtm {
	case FooDTM:
		foo, _, err := FooDTS.UnmarshalData(bs[n:])
		if err != nil {
			panic(err)
		}
		// process foo ...
		fmt.Println(foo)
	case BarDTM:
		bar, _, err := BarDTS.UnmarshalData(bs[n:])
		if err != nil {
			panic(err)
		}
		// process bar ...
		fmt.Println(bar)
	default:
		panic(fmt.Sprintf("unexpected %v DTM", dtm))
	}
}

func randomData() (bs []byte) {
	// Generate random DTM
	dtm := com.DTM(rand.Intn(2) + 1)
	switch dtm {
	// Marshal Foo
	case FooDTM:
		foo := Foo{num: 5}
		bs = make([]byte, FooDTS.Size(foo))
		FooDTS.Marshal(foo, bs)
	// Marshal Bar
	case BarDTM:
		bar := Bar{str: "hello world"}
		bs = make([]byte, BarDTS.Size(bar))
		BarDTS.Marshal(bar, bs)
	default:
		panic(fmt.Sprintf("unexpected %v DTM", dtm))
	}
	return
}
