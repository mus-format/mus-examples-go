package main

import (
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// DTM
// -----------------------------------------------------------------------------

const (
	CopyDTM com.DTM = iota
	InsertDTM
)

// -----------------------------------------------------------------------------
//	Marshal/Unmarshal/Size functions
// -----------------------------------------------------------------------------

// Instruction interface

// With help of the type switch and regular switch we can implement
// Marshal/Unmarshal/Size functions for the Instruction interface.

func MarshalInstructionMUS(instr Instruction, bs []byte) (n int) {
	switch in := instr.(type) {
	case Copy:
		return CopyDTS.Marshal(in, bs)
	case Insert:
		return InsertDTS.Marshal(in, bs)
	default:
		panic(ErrUnexpectedInstructionType)
	}
}

func UnmarshalInstructionMUS(bs []byte) (instr Instruction, n int, err error) {
	dtm, n, err := dts.UnmarshalDTM(bs)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case CopyDTM:
		instr, n1, err = CopyDTS.UnmarshalData(bs[n:])
		n += n1
		return
	case InsertDTM:
		instr, n1, err = InsertDTS.UnmarshalData(bs[n:])
		n += n1
		return
	default:
		err = ErrUnexpectedDTM
		return
	}
}

func SizeInstructionMUS(instr Instruction) (size int) {
	switch in := instr.(type) {
	case Copy:
		return CopyDTS.Size(in)
	case Insert:
		return InsertDTS.Size(in)
	default:
		panic(ErrUnexpectedInstructionType)
	}
}

// Copy

func MarshalCopyMUS(c Copy, bs []byte) (n int) {
	n = varint.MarshalInt(c.start, bs)
	n += varint.MarshalInt(c.end, bs[n:])
	return
}

func UnmarshalCopyMUS(bs []byte) (c Copy, n int, err error) {
	c.start, n, err = varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var n1 int
	c.end, n1, err = varint.UnmarshalInt(bs[n:])
	n += n1
	return
}

func SizeCopyMUS(c Copy) (size int) {
	size = varint.SizeInt(c.start)
	return size + varint.SizeInt(c.end)
}

func SkipCopyMUS(bs []byte) (n int, err error) {
	n, err = varint.SkipInt(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.SkipInt(bs[n:])
	n += n1
	return
}

// Insert

func MarshalInsertMUS(i Insert, bs []byte) (n int) {
	return ord.MarshalString(i.str, nil, bs)
}

func UnmarshalInsertMUS(bs []byte) (i Insert, n int, err error) {
	i.str, n, err = ord.UnmarshalString(nil, bs)
	return
}

func SizeInsertMUS(i Insert) (size int) {
	return ord.SizeString(i.str, nil)
}

func SkipInsertMUS(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}

// -----------------------------------------------------------------------------
// DTS
// -----------------------------------------------------------------------------

var (
	CopyDTS = dts.New[Copy](CopyDTM,
		mus.MarshallerFn[Copy](MarshalCopyMUS),
		mus.UnmarshallerFn[Copy](UnmarshalCopyMUS),
		mus.SizerFn[Copy](SizeCopyMUS),
		mus.SkipperFn(SkipCopyMUS),
	)
	InsertDTS = dts.New[Insert](InsertDTM,
		mus.MarshallerFn[Insert](MarshalInsertMUS),
		mus.UnmarshallerFn[Insert](UnmarshalInsertMUS),
		mus.SizerFn[Insert](SizeInsertMUS),
		mus.SkipperFn(SkipInsertMUS),
	)
)
