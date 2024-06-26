package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/pm"
)

func MarshalTwoPtrMUS(p TwoPtr, bs []byte) (n int) {
	mp := com.NewPtrMap()
	n = pm.MarshalPtr[string](p.ptr1,
		mus.MarshallerFn[string](MarshalString), mp, bs)
	var n1 int
	n1 += pm.MarshalPtr[string](p.ptr2,
		mus.MarshallerFn[string](MarshalString), mp, bs[n:])
	return n + n1
}

func UnmarshalTwoPtrMUS(bs []byte) (p TwoPtr, n int, err error) {
	mp := com.NewReversePtrMap()
	p.ptr1, n, err = pm.UnmarshalPtr[string](
		mus.UnmarshallerFn[string](UnmarshalString), mp, bs)
	if err != nil {
		return
	}
	var n1 int
	p.ptr2, n1, err = pm.UnmarshalPtr[string](
		mus.UnmarshallerFn[string](UnmarshalString), mp, bs[n:])
	n += n1
	return
}

func SizeTowPtrMUS(p TwoPtr) (size int) {
	mp := com.NewPtrMap()
	size += pm.SizePtr[string](p.ptr1, mus.SizerFn[string](SizeString), mp)
	size += pm.SizePtr[string](p.ptr2, mus.SizerFn[string](SizeString), mp)
	return
}

func MarshalString(v string, bs []byte) (n int) {
	return ord.MarshalString(v, nil, bs)
}

func UnmarshalString(bs []byte) (v string, n int, err error) {
	return ord.UnmarshalString(nil, bs)
}

func SizeString(v string) (size int) {
	return ord.SizeString(v, nil)
}
