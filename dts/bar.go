package main

import dts "github.com/mus-format/mus-dts-go"

type Bar struct {
	str string
}

type BarDTSWrap struct {
	DTS dts.DTS[Bar]
}
