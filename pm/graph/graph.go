package main

import (
	"bytes"
	"fmt"
)

func NewGraph[T comparable, V any]() Graph[T, V] {
	return Graph[T, V]{
		Vertices: make(map[T]*Vertex[T, V]),
	}
}

type Graph[T comparable, V any] struct {
	Vertices map[T]*Vertex[T, V]
}

func (g *Graph[T, V]) AddVertex(id T, val V) {
	g.Vertices[id] = &Vertex[T, V]{
		Val: val,
	}
}

func (g *Graph[T, V]) AddEdge(from, to T, weight int) (err error) {
	vFrom, pst := g.Vertices[from]
	if !pst {
		return fmt.Errorf("vertex %v not found", from)
	}
	vTo, pst := g.Vertices[to]
	if !pst {
		return fmt.Errorf("vertex %v not found", from)
	}
	e := &Edge[T, V]{
		Vertex: vTo,
		Weight: weight,
	}
	if vFrom.Edges == nil {
		vFrom.Edges = make(map[T]*Edge[T, V])
	}
	vFrom.Edges[to] = e
	return
}

func (g *Graph[T, V]) String() string {
	var buf bytes.Buffer
	for id, v := range g.Vertices {
		buf.WriteString(fmt.Sprintf("%v: %v\n", id, v.Val))
		for id, e := range v.Edges {
			buf.WriteString(fmt.Sprintf("  %v: %v\n", id, e.Weight))
		}
	}
	return buf.String()
}
