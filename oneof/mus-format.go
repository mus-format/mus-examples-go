package main

import (
	"fmt"
	"reflect"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-go"
	"github.com/mus-format/ext-mus-go"
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
	CopyMUS        = copyMUS{}
	InsertMUS      = insertMUS{}
	InstructionMUS = instructionMUS{}
)

// DTS (Data Type metadata Support) definitions.
var (
	CopyDTS   = dts.New[Copy](CopyDTM, CopyMUS)
	InsertDTS = dts.New[Insert](InsertDTM, InsertMUS)
)

// instructionMUS implements mus.Serializer for Instruction.
type instructionMUS struct{}

func (s instructionMUS) Marshal(in Instruction, bs []byte) (n int) {
	if m, ok := in.(ext.MarshallerTypedMUS); ok {
		return m.MarshalTypedMUS(bs)
	}
	panic(fmt.Sprintf("%v doesn't implement the ext.MarshallerTypedMUS interface", reflect.TypeOf(in)))
}

func (s instructionMUS) Unmarshal(bs []byte) (in Instruction, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	switch dtm {
	case CopyDTM:
		in, n1, err = CopyDTS.UnmarshalData(bs[n:])
		n += n1
	case InsertDTM:
		in, n1, err = InsertDTS.UnmarshalData(bs[n:])
		n += n1
	default:
		err = ErrUnexpectedDTM
	}
	return
}

func (s instructionMUS) Size(in Instruction) (size int) {
	if s, ok := in.(ext.MarshallerTypedMUS); ok {
		return s.SizeTypedMUS()
	}
	panic(fmt.Sprintf("%v doesn't implement the ext.MarshallerTypedMUS interface", reflect.TypeOf(in)))
}

// copyMUS implements mus.Serializer for Copy.
type copyMUS struct{}

func (s copyMUS) Marshal(c Copy, bs []byte) (n int) {
	n = varint.Int.Marshal(c.start, bs)
	n += varint.Int.Marshal(c.end, bs[n:])
	return
}

func (s copyMUS) Unmarshal(bs []byte) (c Copy, n int, err error) {
	c.start, n, err = varint.Int.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	c.end, n1, err = varint.Int.Unmarshal(bs[n:])
	n += n1
	return
}

func (s copyMUS) Size(c Copy) (size int) {
	size = varint.Int.Size(c.start)
	return size + varint.Int.Size(c.end)
}

func (s copyMUS) Skip(bs []byte) (n int, err error) {
	n, err = varint.Int.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = varint.Int.Skip(bs[n:])
	n += n1
	return
}

// insertMUS implements mus.Serializer for Insert.
type insertMUS struct{}

func (s insertMUS) Marshal(i Insert, bs []byte) (n int) {
	return ord.String.Marshal(i.str, bs)
}

func (s insertMUS) Unmarshal(bs []byte) (i Insert, n int, err error) {
	i.str, n, err = ord.String.Unmarshal(bs)
	return
}

func (s insertMUS) Size(i Insert) (size int) {
	return ord.String.Size(i.str)
}

func (s insertMUS) Skip(bs []byte) (n int, err error) {
	return ord.SkipString(nil, bs)
}
