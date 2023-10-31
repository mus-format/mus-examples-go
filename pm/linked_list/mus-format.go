package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Elem
// -----------------------------------------------------------------------------

// NewElemMarshaller creates a new Elem marshaller.
//
// valM param is an Elem value marshaller, elemM - Elem marshaller.
func NewElemMarshaller[T any](valM mus.Marshaller[T],
	eleM *mus.Marshaller[*Elem[T]]) mus.Marshaller[Elem[T]] {
	return mus.MarshallerFn[Elem[T]](
		func(e Elem[T], bs []byte) (n int) {
			n = valM.MarshalMUS(e.Val, bs)
			n += (*eleM).MarshalMUS(e.prev, bs[n:])
			return n + (*eleM).MarshalMUS(e.next, bs[n:])
		},
	)
}

// NewElemUnmarshaller creates a new Elem unmarshaller.
//
// valU param is an Elem value unmarshaller, elemU - Elem unmarshaller.
func NewElemUnmarshaller[T any](valU mus.Unmarshaller[T],
	elemU *mus.Unmarshaller[*Elem[T]]) mus.Unmarshaller[Elem[T]] {
	return mus.UnmarshallerFn[Elem[T]](
		func(bs []byte) (e Elem[T], n int, err error) {
			e.Val, n, err = valU.UnmarshalMUS(bs)
			if err != nil {
				return
			}
			var n1 int
			e.prev, n1, err = (*elemU).UnmarshalMUS(bs[n:])
			n += n1
			if err != nil {
				return
			}
			e.next, n1, err = (*elemU).UnmarshalMUS(bs[n:])
			n += n1
			return
		},
	)
}

// NewElemSizer creates a new Elem sizer.
//
// valS param is an Elem value sizer, elemS - Elem sizer.
func NewElemSizer[T any](valS mus.Sizer[T],
	elemS *mus.Sizer[*Elem[T]]) mus.Sizer[Elem[T]] {
	return mus.SizerFn[Elem[T]](
		func(e Elem[T]) (size int) {
			size = valS.SizeMUS(e.Val)
			size += (*elemS).SizeMUS(e.prev)
			return size + (*elemS).SizeMUS(e.next)
		},
	)
}

// -----------------------------------------------------------------------------
// LinkedList
// -----------------------------------------------------------------------------

// NewLinkedListMarshaller creates a new LinkedList marshaller.
//
// valM param is an Elem value marshaller.
func NewLinkedListMarshaller[T any](
	valM mus.Marshaller[T]) mus.Marshaller[LinkedList[T]] {
	var (
		mp                = com.NewPtrMap()
		elemPtrMarshaller mus.Marshaller[*Elem[T]]
	)
	elemPtrMarshaller = mus.MarshallerFn[*Elem[T]](
		func(v *Elem[T], bs []byte) (n int) {
			return pm.MarshalPtr[Elem[T]](v,
				NewElemMarshaller(valM, &elemPtrMarshaller),
				mp,
				bs)
		},
	)
	return mus.MarshallerFn[LinkedList[T]](
		func(l LinkedList[T], bs []byte) (n int) {
			n = elemPtrMarshaller.MarshalMUS(l.head, bs)
			n += elemPtrMarshaller.MarshalMUS(l.tail, bs[n:])
			return n + varint.MarshalInt(l.len, bs[n:])
		},
	)
}

// NewLinkedListUnmarshaller cretaes a new LinkedList unmarshaller.
//
// valU param is an Elem value unmarshaller.
func NewLinkedListUnmarshaller[T any](
	valU mus.Unmarshaller[T]) mus.Unmarshaller[LinkedList[T]] {
	var (
		mp                  = com.NewReversePtrMap()
		elemPtrUnmarshaller mus.Unmarshaller[*Elem[T]]
	)
	elemPtrUnmarshaller = mus.UnmarshallerFn[*Elem[T]](
		func(bs []byte) (e *Elem[T], n int, err error) {
			return pm.UnmarshalPtr[Elem[T]](
				NewElemUnmarshaller[T](valU, &elemPtrUnmarshaller), mp, bs)
		},
	)
	return mus.UnmarshallerFn[LinkedList[T]](
		func(bs []byte) (l LinkedList[T], n int, err error) {
			l.head, n, err = elemPtrUnmarshaller.UnmarshalMUS(bs)
			if err != nil {
				return
			}
			var n1 int
			l.tail, n1, err = elemPtrUnmarshaller.UnmarshalMUS(bs[n:])
			n += n1
			if err != nil {
				return
			}
			l.len, n1, err = varint.UnmarshalInt(bs[n:])
			n += n1
			return
		},
	)
}

// NewLinkedListSizer creates a new LinkedList sizer.
//
// valS param is an Elem value sizer.
func NewLinkedListSizer[T any](valS mus.Sizer[T]) mus.Sizer[LinkedList[T]] {
	var (
		mp           = com.NewPtrMap()
		elemPtrSizer mus.Sizer[*Elem[T]]
	)
	elemPtrSizer = mus.SizerFn[*Elem[T]](
		func(e *Elem[T]) (size int) {
			return pm.SizePtr[Elem[T]](e, NewElemSizer[T](valS, &elemPtrSizer), mp)
		},
	)
	return mus.SizerFn[LinkedList[T]](
		func(l LinkedList[T]) (size int) {
			size = elemPtrSizer.SizeMUS(l.head)
			size += elemPtrSizer.SizeMUS(l.tail)
			return size + varint.SizeInt(l.Len())
		},
	)
}
