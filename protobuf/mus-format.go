package main

import (
	"fmt"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/unsafe"
	"github.com/mus-format/mus-go/varint"
	"google.golang.org/protobuf/encoding/protowire"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// This file contains all the data related to the MUS format. Here you can find
// serializers for DataV1 and DataV2.
//
// In fact, when implementing Protobuf encoding with mus-go, there is no need
// for types generated by protoc. It is enough to use simple types like this
// one:
//
// 	type DataV1 struct {
// 	 	Str     string
//	  Bool    bool
//    Int32   int32
//    Float64 float64
//    Slice   []int32
//    Time    Timestamp
//  }
//
// The content of this file should be generated.

// Serializers.
var (
	DataV1Protobuf = dataV1Protobuf{}
	DataV2Protobuf = dataV2Protobuf{}

	LenSer       = lenSer{}
	StringSer    = ord.NewStringSerWith(LenSer)
	SliceSer     = sliceSer[int32]{varint.Int32}
	TimestampSer = timestampSer{}
)

// Protobuf tags of the DataV1 and DataV2 fields.
var (
	strFieldTag     = protowire.EncodeTag(1, protowire.BytesType)
	boolFieldTag    = protowire.EncodeTag(2, protowire.VarintType)
	int32FieldTag   = protowire.EncodeTag(3, protowire.Fixed32Type)
	float64FieldTag = protowire.EncodeTag(4, protowire.Fixed64Type)
	sliceFieldTag   = protowire.EncodeTag(5, protowire.BytesType)
	timeFieldTag    = protowire.EncodeTag(6, protowire.BytesType)
)

// dataV1Protobuf implements the mus.Serializer interface for DataV1.
type dataV1Protobuf struct{}

// Marshal marshals data using Protobuf encoding. Actually, there is nothing
// complicated here. For each field (like data.Str or data.Bool) it:
// 1. Marshals the tag.
// 2. Marshals the length of the value, if the field is a struct or slice.
// 3. Marshals the value.
func (s dataV1Protobuf) Marshal(data *DataV1, bs []byte) (n int) {
	if data.Str != "" {
		n += varint.Uint64.Marshal(strFieldTag, bs[n:])
		n += StringSer.Marshal(data.Str, bs[n:])
	}
	if data.Bool {
		n += varint.Uint64.Marshal(boolFieldTag, bs[n:])
		n += unsafe.Bool.Marshal(data.Bool, bs[n:])
	}
	if data.Int32 != 0 {
		n += varint.Uint64.Marshal(int32FieldTag, bs[n:])
		n += unsafe.Int32.Marshal(data.Int32, bs[n:])
	}
	if data.Float64 != 0 {
		n += varint.Uint64.Marshal(float64FieldTag, bs[n:])
		n += unsafe.Float64.Marshal(data.Float64, bs[n:])
	}
	if len(data.Slice) > 0 {
		n += varint.Uint64.Marshal(sliceFieldTag, bs[n:])
		n += SliceSer.Marshal(data.Slice, bs[n:])
	}
	if data.Time != nil && (data.Time.Seconds != 0 || data.Time.Nanos != 0) {
		n += varint.Uint64.Marshal(timeFieldTag, bs[n:])
		n += TimestampSer.Marshal(data.Time, bs[n:])
	}
	return
}

// UnmarshalDataV1Protobuf unmarshals fields in а loop:
// 1. Unmarshals the tag.
// 2. Unmarshals the length of the value, if the field is a struct or slice.
// 3. Unmarshals the value.
func (s dataV1Protobuf) Unmarshal(bs []byte) (data *DataV1, n int, err error) {
	var (
		n1  int
		l   = len(bs)
		tag uint64
	)
	data = &DataV1{}
	for n < l {
		tag, n1, err = varint.Uint64.Unmarshal(bs[n:])
		n += n1
		if err != nil {
			return
		}
		switch tag {
		case strFieldTag:
			data.Str, n1, err = StringSer.Unmarshal(bs[n:])
		case boolFieldTag:
			data.Bool, n1, err = unsafe.Bool.Unmarshal(bs[n:])
		case int32FieldTag:
			data.Int32, n1, err = unsafe.Int32.Unmarshal(bs[n:])
		case float64FieldTag:
			data.Float64, n1, err = unsafe.Float64.Unmarshal(bs[n:])
		case sliceFieldTag:
			data.Slice, n1, err = SliceSer.Unmarshal(bs[n:])
		case timeFieldTag:
			data.Time, n1, err = TimestampSer.Unmarshal(bs[n:])
		default:
			err = fmt.Errorf("unexpected tag %v", tag)
		}
		n += n1
		if err != nil {
			return
		}
	}
	return
}

func (s dataV1Protobuf) Size(data *DataV1) (size int) {
	if data.Str != "" {
		size += varint.Uint64.Size(strFieldTag)
		size += StringSer.Size(data.Str)
	}
	if data.Bool {
		size += varint.Uint64.Size(boolFieldTag)
		size += unsafe.Bool.Size(data.Bool)
	}
	if data.Int32 != 0 {
		size += varint.Uint64.Size(int32FieldTag)
		size += unsafe.Int32.Size(data.Int32)
	}
	if data.Float64 != 0 {
		size += varint.Uint64.Size(float64FieldTag)
		size += unsafe.Float64.Size(data.Float64)
	}
	if len(data.Slice) > 0 {
		size += varint.Uint64.Size(sliceFieldTag)
		size += SliceSer.Size(data.Slice)
	}
	if data.Time != nil && (data.Time.Seconds != 0 || data.Time.Nanos != 0) {
		size += varint.Uint64.Size(timeFieldTag)
		size += TimestampSer.Size(data.Time)
	}
	return
}

// dataV2Protobuf implements the mus.Serializer interface for DataV2.
type dataV2Protobuf struct{}

func (s dataV2Protobuf) Marshal(data *DataV2, bs []byte) (n int) {
	if data.Str != "" {
		n += varint.Uint64.Marshal(strFieldTag, bs[n:])
		n += StringSer.Marshal(data.Str, bs[n:])
	}
	if data.Int32 != 0 {
		n += varint.Uint64.Marshal(int32FieldTag, bs[n:])
		n += unsafe.Int32.Marshal(data.Int32, bs[n:])
	}
	if data.Float64 != 0 {
		n += varint.Uint64.Marshal(float64FieldTag, bs[n:])
		n += unsafe.Float64.Marshal(data.Float64, bs[n:])
	}
	if data.Time != nil && (data.Time.Seconds != 0 || data.Time.Nanos != 0) {
		n += varint.Uint64.Marshal(timeFieldTag, bs[n:])
		n += TimestampSer.Marshal(data.Time, bs[n:])
	}
	return
}

func (s dataV2Protobuf) Unmarshal(bs []byte) (data *DataV2, n int, err error) {
	var (
		n1  int
		l   = len(bs)
		tag uint64
	)
	data = &DataV2{}
	for n < l {
		tag, n1, err = varint.Uint64.Unmarshal(bs[n:])
		n += n1
		if err != nil {
			return
		}
		switch tag {
		case strFieldTag:
			data.Str, n1, err = StringSer.Unmarshal(bs[n:])
		case boolFieldTag:
			// Bool field was removed in DataV2, so simply skip it here.
			n1, err = unsafe.Bool.Skip(bs[n:])
		case int32FieldTag:
			data.Int32, n1, err = unsafe.Int32.Unmarshal(bs[n:])
		case float64FieldTag:
			data.Float64, n1, err = unsafe.Float64.Unmarshal(bs[n:])
		case sliceFieldTag:
			// Slice field was remove in DataV2, so simply skip it here.
			n1, err = SliceSer.Skip(bs[n:])
		case timeFieldTag:
			data.Time, n1, err = TimestampSer.Unmarshal(bs[n:])
		default:
			err = fmt.Errorf("unexpected tag %v", tag)
		}
		n += n1
		if err != nil {
			return
		}
	}
	return
}

func (s dataV2Protobuf) Size(data *DataV2) (size int) {
	if data.Str != "" {
		size += varint.Uint64.Size(strFieldTag)
		size += StringSer.Size(data.Str)
	}
	if data.Int32 != 0 {
		size += varint.Uint64.Size(int32FieldTag)
		size += unsafe.Int32.Size(data.Int32)
	}
	if data.Float64 != 0 {
		size += varint.Uint64.Size(float64FieldTag)
		size += unsafe.Float64.Size(data.Float64)
	}
	if data.Time != nil && (data.Time.Seconds != 0 || data.Time.Nanos != 0) {
		size += varint.Uint64.Size(timeFieldTag)
		size += TimestampSer.Size(data.Time)
	}
	return
}

var (
	secondsFieldTag = protowire.EncodeTag(1, protowire.VarintType)
	nanosFieldTag   = protowire.EncodeTag(2, protowire.VarintType)
)

// timestampSer implements the mus.Serializer interface for timestamppb.Timestamp.
type timestampSer struct{}

func (s timestampSer) Marshal(tm *timestamppb.Timestamp, bs []byte) (n int) {
	size := s.size(tm)
	if size > 0 {
		n += varint.PositiveInt.Marshal(size, bs[n:])
		if tm.Seconds != 0 {
			n += varint.Uint64.Marshal(secondsFieldTag, bs[n:])
			n += varint.PositiveInt64.Marshal(tm.Seconds, bs[n:])
		}
		if tm.Nanos != 0 {
			n += varint.Uint64.Marshal(nanosFieldTag, bs[n:])
			n += varint.PositiveInt32.Marshal(tm.Nanos, bs[n:])
		}
	}
	return
}

func (timestampSer) Unmarshal(bs []byte) (tm *timestamppb.Timestamp, n int,
	err error) {
	n, err = varint.PositiveInt.Skip(bs)
	if err != nil {
		return
	}
	var (
		n1  int
		l   = len(bs)
		tag uint64
	)
	tm = &timestamppb.Timestamp{}
	for {
		tag, n1, err = varint.Uint64.Unmarshal(bs[n:])
		n += n1
		if err != nil {
			return
		}
		switch tag {
		case secondsFieldTag:
			tm.Seconds, n1, err = varint.PositiveInt64.Unmarshal(bs[n:])
		case nanosFieldTag:
			tm.Nanos, n1, err = varint.PositiveInt32.Unmarshal(bs[n:])
		}
		n += n1
		if err != nil {
			return
		}
		if n == l {
			return
		}
	}
}

func (s timestampSer) Size(tm *timestamppb.Timestamp) (size int) {
	size = s.size(tm)
	return size + varint.PositiveInt.Size(size)
}

func (s timestampSer) size(tm *timestamppb.Timestamp) (size int) {
	if tm.Seconds != 0 {
		size += varint.Uint64.Size(secondsFieldTag)
		size += varint.PositiveInt64.Size(tm.Seconds)
	}
	if tm.Nanos != 0 {
		size += varint.Uint64.Size(nanosFieldTag)
		size += varint.PositiveInt32.Size(tm.Nanos)
	}
	return
}

// sliceSer implements the mus.Serializer interface for slices.
type sliceSer[T any] struct {
	elemSer mus.Serializer[T]
}

func (s sliceSer[T]) Marshal(sl []T, bs []byte) (n int) {
	length := len(sl)
	if length > 0 {
		n += varint.PositiveInt.Marshal(s.size(sl), bs[n:])
		for i := 0; i < len(sl); i++ {
			n += s.elemSer.Marshal(sl[i], bs[n:])
		}
	}
	return
}

func (s sliceSer[T]) Unmarshal(bs []byte) (sl []T, n int, err error) {
	var (
		n1   int
		elem T
	)
	sl = []T{}
	size, n, err := varint.PositiveInt.Unmarshal(bs)
	if err != nil {
		return
	}
	if len(bs) < size {
		err = com.ErrOverflow
		return
	}
	for n < size {
		elem, n1, err = s.elemSer.Unmarshal(bs[n:])
		n += n1
		if err != nil {
			return
		}
		sl = append(sl, elem)
	}
	return
}

func (s sliceSer[T]) Size(sl []T) (size int) {
	size = s.size(sl)
	return size + varint.PositiveInt.Size(size)
}

func (s sliceSer[T]) Skip(bs []byte) (n int, err error) {
	l, n, err := varint.PositiveInt.Unmarshal(bs)
	if err != nil {
		return
	}
	n += l
	return
}

func (s sliceSer[T]) size(sl []T) (size int) {
	for i := 0; i < len(sl); i++ {
		size += s.elemSer.Size(sl[i])
	}
	return
}

// lenSer implements the mus.Serializer interface for length.
type lenSer struct{}

func (lenSer) Marshal(v int, bs []byte) (n int) {
	return varint.PositiveInt32.Marshal(int32(v), bs)
}

func (lenSer) Unmarshal(bs []byte) (v int, n int, err error) {
	v32, n, err := varint.PositiveInt32.Unmarshal(bs)
	v = int(v32)
	return
}

func (lenSer) Size(v int) (size int) {
	return varint.PositiveInt32.Size(int32(v))
}

func (lenSer) Skip(bs []byte) (n int, err error) {
	return varint.PositiveInt32.Skip(bs)
}
