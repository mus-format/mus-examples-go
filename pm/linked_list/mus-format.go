package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// MakeLinkedListMUS creates a new LinkedList serializer that can be directly used.
func MakeLinkedListMUS[T any](valMUS mus.Serializer[T]) mus.Serializer[LinkedList[T]] {
	var (
		ptrMap    = com.NewPtrMap()
		revPtrMap = com.NewReversePtrMap()
		ser       = NewLinkedListMUS(ptrMap, revPtrMap, valMUS)
	)
	return pm.Wrap[LinkedList[T]](ptrMap, revPtrMap, ser)
}

// NewLinkedListMUS creates a new LinkedList serializer.
func NewLinkedListMUS[T any](ptrMap *com.PtrMap, revPtrMap *com.ReversePtrMap,
	valMUS mus.Serializer[T],
) linkedListMUS[T] {
	var (
		elPtrMUS mus.Serializer[*Elem[T]]
		elMUS    = elemMUS[T]{valMUS, &elPtrMUS}
	)
	elPtrMUS = pm.NewPtrSer[Elem[T]](ptrMap, revPtrMap, elMUS)
	return linkedListMUS[T]{elPtrMUS}
}

// elemMUS implements mus.Serializer for Elem.
type elemMUS[T any] struct {
	valMUS mus.Serializer[T]
	elMUS  *mus.Serializer[*Elem[T]]
}

func (s elemMUS[T]) Marshal(e Elem[T], bs []byte) (n int) {
	n = s.valMUS.Marshal(e.Val, bs)
	n += (*s.elMUS).Marshal(e.prev, bs[n:])
	return n + (*s.elMUS).Marshal(e.next, bs[n:])
}

func (s elemMUS[T]) Unmarshal(bs []byte) (e Elem[T], n int, err error) {
	e.Val, n, err = s.valMUS.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	e.prev, n1, err = (*s.elMUS).Unmarshal(bs[n:])
	n += n1
	if err != nil {
		return
	}
	e.next, n1, err = (*s.elMUS).Unmarshal(bs[n:])
	n += n1
	return
}

func (s elemMUS[T]) Size(e Elem[T]) (size int) {
	size = s.valMUS.Size(e.Val)
	size += (*s.elMUS).Size(e.prev)
	return size + (*s.elMUS).Size(e.next)
}

func (s elemMUS[T]) Skip(bs []byte) (n int, err error) {
	n, err = s.valMUS.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.elMUS).Skip(bs[n:])
	n += n1
	if err != nil {
		return
	}
	n1, err = (*s.elMUS).Skip(bs[n:])
	n += n1
	return
}

// linkedListMUS implements mus.Serializer for LinkedList.
type linkedListMUS[T any] struct {
	elPtrMUS mus.Serializer[*Elem[T]]
}

func (s linkedListMUS[T]) Marshal(l LinkedList[T], bs []byte) (n int) {
	n = s.elPtrMUS.Marshal(l.head, bs)
	n += s.elPtrMUS.Marshal(l.tail, bs[n:])
	return n + varint.PositiveInt.Marshal(l.len, bs[n:])
}

func (s linkedListMUS[T]) Unmarshal(bs []byte) (l LinkedList[T], n int, err error) {
	l.head, n, err = s.elPtrMUS.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	l.tail, n1, err = s.elPtrMUS.Unmarshal(bs[n:])
	n += n1
	if err != nil {
		return
	}
	l.len, n1, err = varint.PositiveInt.Unmarshal(bs[n:])
	n += n1
	return
}

func (s linkedListMUS[T]) Size(l LinkedList[T]) (size int) {
	size = s.elPtrMUS.Size(l.head)
	size += s.elPtrMUS.Size(l.tail)
	return size + varint.PositiveInt.Size(l.len)
}

func (s linkedListMUS[T]) Skip(bs []byte) (n int, err error) {
	n, err = s.elPtrMUS.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = s.elPtrMUS.Skip(bs[n:])
	n += n1
	if err != nil {
		return
	}
	n1, err = varint.PositiveInt.Skip(bs[n:])
	n += n1
	return
}
