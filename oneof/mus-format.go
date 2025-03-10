package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// DTM (Data Type Metadata) definitions.
const (
	CopyDTM com.DTM = iota
	InsertDTM
)

// Serializers.
var (
	CopySer        = copySer{}
	InsertSer      = insertSer{}
	InstructionSer = instructionSer{}
)

// DTS (Data Type Metadata Support) definitions.
var (
	CopyDTS   = dts.New[Copy](CopyDTM, CopySer)
	InsertDTS = dts.New[Insert](InsertDTM, InsertSer)
)

// instructionSer implements mus.Serializer for Instruction.
type instructionSer struct{}

func (s instructionSer) Marshal(in Instruction, bs []byte) (n int) {
	if m, ok := in.(MarshallerMUS); ok {
		return m.MarshalMUS(bs)
	}
	panic("in doesn't implement MarshallerMUS interface")
}

func (s instructionSer) Unmarshal(bs []byte) (in Instruction, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case CopyDTM:
		in, n1, err = CopyDTS.UnmarshalData(bs[n:])
		n += n1
		return
	case InsertDTM:
		in, n1, err = InsertDTS.UnmarshalData(bs[n:])
		n += n1
		return
	default:
		err = ErrUnexpectedDTM
		return
	}
}

func (s instructionSer) Size(in Instruction) (size int) {
	if s, ok := in.(MarshallerMUS); ok {
		return s.SizeMUS()
	}
	panic("in doesn't implement MarshallerMUS interface")
}

// copySer implements mus.Serializer for Copy.
type copySer struct{}

func (s copySer) Marshal(c Copy, bs []byte) (n int) {
	n = varint.Int.Marshal(c.start, bs)
	n += varint.Int.Marshal(c.end, bs[n:])
	return
}

func (s copySer) Unmarshal(bs []byte) (c Copy, n int, err error) {
	c.start, n, err = varint.Int.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	c.end, n1, err = varint.Int.Unmarshal(bs[n:])
	n += n1
	return
}

func (s copySer) Size(c Copy) (size int) {
	size = varint.Int.Size(c.start)
	return size + varint.Int.Size(c.end)
}

func (s copySer) Skip(bs []byte) (n int, err error) {
	n, err = varint.Int.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.Int.Skip(bs[n:])
	n += n1
	return
}

// insertSer implements mus.Serializer for Insert.
type insertSer struct{}

func (s insertSer) Marshal(i Insert, bs []byte) (n int) {
	return ord.String.Marshal(i.str, bs)
}

func (s insertSer) Unmarshal(bs []byte) (i Insert, n int, err error) {
	i.str, n, err = ord.String.Unmarshal(bs)
	return
}

func (s insertSer) Size(i Insert) (size int) {
	return ord.String.Size(i.str)
}

func (s insertSer) Skip(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}
