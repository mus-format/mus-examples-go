package main

import (
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/varint"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

// Demonstrates how to use the pm package to serialize a linked list.
func main() {
	var (
		m = NewLinkedListMarshaller[int](mus.MarshallerFn[int](varint.MarshalInt))
		u = NewLinkedListUnmarshaller[int](mus.UnmarshallerFn[int](varint.UnmarshalInt))
		s = NewLinkedListSizer[int](mus.SizerFn[int](varint.SizeInt))
	)
	l := makeLinkedList()

	// Marshal list.
	bs := make([]byte, s.Size(l))
	m.Marshal(l, bs)

	// Unmarshal list.
	al, _, err := u.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(l, al)
}

func makeLinkedList() (l LinkedList[int]) {
	l = LinkedList[int]{}
	l.AddBack(8)
	l.AddBack(9)
	l.AddBack(10)
	l.AddBack(11)
	return
}
