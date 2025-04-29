package main

import (
	"github.com/mus-format/mus-go/varint"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example demonstrates how to use the pm package to serialize a linked
// list.
func main() {
	var (
		v   = ShortLinkedList()
		ser = MakeLinkedListMUS[int](varint.PositiveInt)
	)

	// Marshal list.
	bs := make([]byte, ser.Size(v))
	ser.Marshal(v, bs)

	// Unmarshal list.
	av, _, err := ser.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(v, av)
}

func ShortLinkedList() (l LinkedList[int]) {
	l = LinkedList[int]{}
	l.AddBack(8)
	l.AddBack(9)
	l.AddBack(10)
	l.AddBack(11)
	return
}
