package main

import "github.com/mus-format/mus-go/ord"

type Bar struct {
	a string
	b bool
}

// -----------------------------------------------------------------------------
func MarshalMetaBar(bar Bar, bs []byte) (n int) {
	n = MarshalDataType(BarType, bs)
	return n + MarshalBar(bar, bs[n:])
}

func SizeMetaBar(bar Bar) (size int) {
	size += SizeDataType(BarType)
	return size + SizeBar(bar)
}

// -----------------------------------------------------------------------------
func MarshalBar(bar Bar, bs []byte) (n int) {
	n = ord.MarshalString(bar.a, bs)
	return n + ord.MarshalBool(bar.b, bs[n:])
}

func UnmarshalBar(bs []byte) (bar Bar, n int, err error) {
	bar.a, n, err = ord.UnmarshalString(bs)
	if err != nil {
		return
	}
	var n1 int
	bar.b, n1, err = ord.UnmarshalBool(bs)
	n += n1
	return
}

func SizeBar(bar Bar) (size int) {
	size += ord.SizeString(bar.a)
	return size + ord.SizeBool(bar.b)
}
