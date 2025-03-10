package main

import (
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
	assert "github.com/ymz-ncnk/assert/panic"
)

func init() {
	assert.On = true
}

// This example demonstrates how to use the pm package to serialize a cyclic
// graph.
func main() {
	var (
		v   = CyclicGraph()
		ser = MakeGraphSer[int, string](varint.Int, ord.String)
	)

	// Marshal graph.
	bs := make([]byte, ser.Size(v))
	ser.Marshal(v, bs)

	// Unmarshal graph.
	av, _, err := ser.Unmarshal(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(v, av)
}

func CyclicGraph() (g Graph[int, string]) {
	g = NewGraph[int, string]()
	g.AddVertex(1, "one")
	g.AddVertex(2, "two")
	g.AddVertex(3, "three")

	g.AddEdge(1, 2, 10)
	g.AddEdge(2, 3, 20)
	g.AddEdge(3, 1, 30)
	return
}
