package main

import dts "github.com/mus-format/mus-dts-go"

type Foo struct {
	num int
}

type FooDTSWrap struct {
	DTS dts.DTS[Foo]
}

func (w FooDTSWrap) MarshalMUS() {

}
