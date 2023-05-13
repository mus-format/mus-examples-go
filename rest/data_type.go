package main

import "github.com/mus-format/mus-go/raw"

// DataType in this example defines a product version.
type DataType byte

func MarshalDataType(dtm DataType, bs []byte) (n int) {
	return raw.MarshalByte(byte(dtm), bs)
}

func UnmarshalDataType(bs []byte) (dtm DataType, n int,
	err error) {
	b, n, err := raw.UnmarshalByte(bs)
	dtm = DataType(b)
	return
}

func SizeDataType(dtm DataType) (size int) {
	return raw.SizeByte(byte(dtm))
}
