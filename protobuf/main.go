package main

import (
	"errors"

	"github.com/brianvoe/gofakeit"
	"github.com/ymz-ncnk/assert"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	assert.On = true
}

// In this example we use mus-go to implement Protobuf encoding.
func main() {
	// We have two structures DataV1, DataV2. The last one is the same as DataV1,
	// but with two deleted fields: Bool and Slice. For both of these structures,
	// we:
	// - generate ...pb.go files using protoc.
	// - write mus-format.go (actually, this code could be generated, but there is
	//   no such generator yet).
	var (
		dataV1 = DataV1{
			Str:     gofakeit.UUID(),
			Bool:    gofakeit.Bool(),
			Int32:   gofakeit.Int32(),
			Float64: gofakeit.Float64(),
			Slice:   []int32{gofakeit.Int32(), gofakeit.Int32()},
			Time:    timestamppb.New(gofakeit.Date()),
		}
		dataV2 = DataV2{
			Str:     gofakeit.UUID(),
			Int32:   gofakeit.Int32(),
			Float64: gofakeit.Float64(),
			Time:    timestamppb.New(gofakeit.Date()),
		}
	)
	// Let's marshal using protobuf and Unmarshal using mus-go implementation (the
	// Unmarshalled data is compared with the original at the end).
	MarshalProtobuf_UnmarshalMusGo(&dataV1)
	// Marshal using mus-go - Unmarshal using protobuf.
	MarshalMusGo_UnmarshalProtobuf(&dataV1)
	// Marshal first version and Unmarshal second, both using mus-go.
	MarshalDataV1_UnmarshalDataV2(&dataV1)
	// Marshal second version and Unmarshal first one again using mus-go.
	MarshalDataV2_UnmarshalDataV1(&dataV2)

	// As you can see, everything works as expected.
}

func MarshalProtobuf_UnmarshalMusGo(data *DataV1) {
	bs, err := proto.Marshal(data)
	assert.EqualError(err, nil)

	adata, _, err := UnmarshalDataV1Protobuf(bs)
	assert.EqualError(err, nil)

	assert.Equal(data.String(), adata.String())
}

func MarshalMusGo_UnmarshalProtobuf(data *DataV1) {
	bs := make([]byte, SizeDataV1Protobuf(data))
	MarshalDataV1Protobuf(data, bs)

	adata := DataV1{}
	err := proto.Unmarshal(bs, &adata)
	assert.EqualError(err, nil)

	assert.Equal(data.String(), adata.String())
}

func MarshalDataV1_UnmarshalDataV2(dataV1 *DataV1) {
	bs := make([]byte, SizeDataV1Protobuf(dataV1))
	MarshalDataV1Protobuf(dataV1, bs)

	dataV2, _, err := UnmarshalDataV2Protobuf(bs)
	assert.EqualError(err, nil)

	if err := same(dataV1, dataV2); err != nil {
		panic(err)
	}
}

func MarshalDataV2_UnmarshalDataV1(dataV2 *DataV2) {
	bs := make([]byte, SizeDataV2Protobuf(dataV2))
	MarshalDataV2Protobuf(dataV2, bs)

	dataV1, _, err := UnmarshalDataV1Protobuf(bs)
	assert.EqualError(err, nil)

	if err := same(dataV1, dataV2); err != nil {
		panic(err)
	}
}

func same(dataV1 *DataV1, dataV2 *DataV2) (err error) {
	if dataV1.Str != dataV2.Str {
		return errors.New("Str")
	}
	if dataV1.Int32 != dataV2.Int32 {
		return errors.New("Int32")
	}
	if dataV1.Float64 != dataV2.Float64 {
		return errors.New("Float64")
	}
	if dataV1.Time != nil && dataV2.Time != nil {
		if dataV1.Time.Seconds != dataV2.Time.Seconds {
			return errors.New("Seconds")
		}
		if dataV1.Time.Nanos != dataV2.Time.Nanos {
			return errors.New("Nanos")
		}
	}
	if (dataV1.Time != nil && dataV2.Time == nil) ||
		(dataV1.Time == nil && dataV2.Time != nil) {
		return errors.New("one time is nil, another is not")
	}
	return
}
