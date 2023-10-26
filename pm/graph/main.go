package main

import (
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/varint"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

func main() {
	var (
		m = NewGraphMarshaller[int, int](mus.MarshallerFn[int](varint.MarshalInt),
			mus.MarshallerFn[int](varint.MarshalInt))
		u = NewGraphUnmarshaller[int, int](mus.UnmarshallerFn[int](varint.UnmarshalInt),
			mus.UnmarshallerFn[int](varint.UnmarshalInt))
		s = NewGraphSizer[int, int](mus.SizerFn[int](varint.SizeInt),
			mus.SizerFn[int](varint.SizeInt))
		g = makeCyclicGraph()
	)
	bs := make([]byte, s.SizeMUS(g))
	m.MarshalMUS(g, bs)
	ag, _, err := u.UnmarshalMUS(bs)
	assert.EqualError(err, nil)
	assert.EqualDeep(g, ag)
}

func makeCyclicGraph() Graph[int, int] {
	var (
		e1, e2, e3 Edge[int]
		v1, v2, v3 Vertex[int]
	)
	e1 = Edge[int]{Weight: 20, Vertex: &v1}
	e2 = Edge[int]{Weight: 10, Vertex: &v2}
	e3 = Edge[int]{Weight: 30, Vertex: &v2}
	v1 = Vertex[int]{Val: 8, Edges: map[int]*Edge[int]{2: &e2}}
	v2 = Vertex[int]{Val: 9, Edges: map[int]*Edge[int]{1: &e1}}
	v3 = Vertex[int]{Val: 10, Edges: map[int]*Edge[int]{2: &e3}}
	return Graph[int, int]{
		Vertices: map[int]*Vertex[int]{1: &v1, 2: &v2, 3: &v3},
	}
}
