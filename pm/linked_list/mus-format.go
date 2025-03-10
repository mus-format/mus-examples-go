package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// MakeLinkedListSer creates a new LinkedList serializer that can be directly used.
func MakeLinkedListSer[T any](valSer mus.Serializer[T]) mus.Serializer[LinkedList[T]] {
	var (
		ptrMap    = com.NewPtrMap()
		revPtrMap = com.NewReversePtrMap()
		ser       = NewLinkedListSer(ptrMap, revPtrMap, valSer)
	)
	return pm.Wrap[LinkedList[T]](ptrMap, revPtrMap, ser)
}

// NewLinkedListSer creates a new LinkedList serializer.
func NewLinkedListSer[T any](ptrMap *com.PtrMap, revPtrMap *com.ReversePtrMap,
	valSer mus.Serializer[T],
) linkedListSer[T] {
	var (
		elPtrSer mus.Serializer[*Elem[T]]
		elSer    = elemSer[T]{valSer, &elPtrSer}
	)
	elPtrSer = pm.NewPtrSer[Elem[T]](ptrMap, revPtrMap, elSer)
	return linkedListSer[T]{elPtrSer}
}

// elemSer implements mus.Serializer for Elem.
type elemSer[T any] struct {
	valSer mus.Serializer[T]
	elSer  *mus.Serializer[*Elem[T]]
}

func (s elemSer[T]) Marshal(e Elem[T], bs []byte) (n int) {
	n = s.valSer.Marshal(e.Val, bs)
	n += (*s.elSer).Marshal(e.prev, bs[n:])
	return n + (*s.elSer).Marshal(e.next, bs[n:])
}

func (s elemSer[T]) Unmarshal(bs []byte) (e Elem[T], n int, err error) {
	e.Val, n, err = s.valSer.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	e.prev, n1, err = (*s.elSer).Unmarshal(bs[n:])
	n += n1
	if err != nil {
		return
	}
	e.next, n1, err = (*s.elSer).Unmarshal(bs[n:])
	n += n1
	return
}

func (s elemSer[T]) Size(e Elem[T]) (size int) {
	size = s.valSer.Size(e.Val)
	size += (*s.elSer).Size(e.prev)
	return size + (*s.elSer).Size(e.next)
}

func (s elemSer[T]) Skip(bs []byte) (n int, err error) {
	n, err = s.valSer.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.elSer).Skip(bs[n:])
	n += n1
	if err != nil {
		return
	}
	n1, err = (*s.elSer).Skip(bs[n:])
	n += n1
	return
}

// linkedListSer implements mus.Serializer for LinkedList.
type linkedListSer[T any] struct {
	elPtrSer mus.Serializer[*Elem[T]]
}

func (s linkedListSer[T]) Marshal(l LinkedList[T], bs []byte) (n int) {
	n = s.elPtrSer.Marshal(l.head, bs)
	n += s.elPtrSer.Marshal(l.tail, bs[n:])
	return n + varint.PositiveInt.Marshal(l.len, bs[n:])
}

func (s linkedListSer[T]) Unmarshal(bs []byte) (l LinkedList[T], n int, err error) {
	l.head, n, err = s.elPtrSer.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	l.tail, n1, err = s.elPtrSer.Unmarshal(bs[n:])
	n += n1
	if err != nil {
		return
	}
	l.len, n1, err = varint.PositiveInt.Unmarshal(bs[n:])
	n += n1
	return
}

func (s linkedListSer[T]) Size(l LinkedList[T]) (size int) {
	size = s.elPtrSer.Size(l.head)
	size += s.elPtrSer.Size(l.tail)
	return size + varint.PositiveInt.Size(l.len)
}

func (s linkedListSer[T]) Skip(bs []byte) (n int, err error) {
	n, err = s.elPtrSer.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = s.elPtrSer.Skip(bs[n:])
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.PositiveInt.Skip(bs[n:])
	n += n1
	return
}
