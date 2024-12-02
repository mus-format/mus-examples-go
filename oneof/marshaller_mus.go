package main

type MarshallerMUS interface {
	MarshalMUS(bs []byte) (n int)
	SizeMUS() (size int)
}
