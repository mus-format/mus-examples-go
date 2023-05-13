package main

import "github.com/mus-format/mus-go/raw"

// Data type metadata.
type DataType byte

func MarshalDataType(dt DataType, bs []byte) (n int) {
	return raw.MarshalByte(byte(dt), bs)
}

func UnmarshalDataType(bs []byte) (dt DataType, n int,
	err error) {
	b, n, err := raw.UnmarshalByte(bs)
	dt = DataType(b)
	return
}

func SizeDataType(dt DataType) (size int) {
	return raw.SizeByte(byte(dt))
}
