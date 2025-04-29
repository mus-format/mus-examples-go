package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// MakeGraphMUS creates a new Graph serializer that can be directly used.
func MakeGraphMUS[T comparable, V any](keyMUS mus.Serializer[T],
	valMUS mus.Serializer[V],
) mus.Serializer[Graph[T, V]] {
	var (
		ptrMap    = com.NewPtrMap()
		revPtrMap = com.NewReversePtrMap()
	)
	return pm.Wrap[Graph[T, V]](ptrMap, revPtrMap, NewGraphMUS[T, V](ptrMap,
		revPtrMap, keyMUS, valMUS))
}

// NewGraphMUS creates a new Graph serializer.
func NewGraphMUS[T comparable, V any](ptrMap *com.PtrMap,
	revPtrMap *com.ReversePtrMap,
	keyMUS mus.Serializer[T],
	valMUS mus.Serializer[V],
) graphMUS[T, V] {
	var (
		edgesMUS    mus.Serializer[map[T]*Edge[T, V]]
		vertecesMUS mus.Serializer[map[T]*Vertex[T, V]]

		edgePtrMUS   mus.Serializer[*Edge[T, V]]
		vertexPtrMUS mus.Serializer[*Vertex[T, V]]

		edgeMUS   = edgeMUS[T, V]{&vertexPtrMUS}
		vertexMUS = vertexMUS[T, V]{valMUS, &edgesMUS}
	)
	edgePtrMUS = pm.NewPtrSer[Edge[T, V]](ptrMap, revPtrMap, edgeMUS)
	vertexPtrMUS = pm.NewPtrSer[Vertex[T, V]](ptrMap, revPtrMap, vertexMUS)

	edgesMUS = ord.NewMapSer[T, *Edge[T, V]](keyMUS, edgePtrMUS)
	vertecesMUS = ord.NewMapSer[T, *Vertex[T, V]](keyMUS, vertexPtrMUS)

	return graphMUS[T, V]{vertecesMUS}
}

// edgeMUS implements the mus.Serializer interface for Edge.
type edgeMUS[T comparable, V any] struct {
	vertexPtrMUS *mus.Serializer[*Vertex[T, V]]
}

func (s edgeMUS[T, V]) Marshal(e Edge[T, V], bs []byte) (n int) {
	n = varint.Int.Marshal(e.Weight, bs)
	return n + (*s.vertexPtrMUS).Marshal(e.Vertex, bs[n:])
}

func (s edgeMUS[T, V]) Unmarshal(bs []byte) (e Edge[T, V], n int, err error) {
	e.Weight, n, err = varint.Int.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	e.Vertex, n1, err = (*s.vertexPtrMUS).Unmarshal(bs[n:])
	n += n1
	return
}

func (s edgeMUS[T, V]) Size(e Edge[T, V]) (size int) {
	size = varint.Int.Size(e.Weight)
	return size + (*s.vertexPtrMUS).Size(e.Vertex)
}

func (s edgeMUS[T, V]) Skip(bs []byte) (n int, err error) {
	n, err = varint.Int.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.vertexPtrMUS).Skip(bs[n:])
	n += n1
	return
}

// vertexMUS implements the mus.Serializer interface for Vertex.
type vertexMUS[T comparable, V any] struct {
	valMUS   mus.Serializer[V]
	edgesMUS *mus.Serializer[map[T]*Edge[T, V]]
}

func (s vertexMUS[T, V]) Marshal(v Vertex[T, V], bs []byte) (n int) {
	n = s.valMUS.Marshal(v.Val, bs)
	return n + (*s.edgesMUS).Marshal(v.Edges, bs[n:])
}

func (s vertexMUS[T, V]) Unmarshal(bs []byte) (v Vertex[T, V], n int, err error) {
	v.Val, n, err = s.valMUS.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	v.Edges, n1, err = (*s.edgesMUS).Unmarshal(bs[n:])
	n += n1
	return
}

func (s vertexMUS[T, V]) Size(v Vertex[T, V]) (size int) {
	size = s.valMUS.Size(v.Val)
	return size + (*s.edgesMUS).Size(v.Edges)
}

func (s vertexMUS[T, V]) Skip(bs []byte) (n int, err error) {
	n, err = s.valMUS.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.edgesMUS).Skip(bs[n:])
	n += n1
	return
}

// graphMUS implements the mus.Serializer interface for Graph.
type graphMUS[T comparable, V any] struct {
	vertecesMUS mus.Serializer[map[T]*Vertex[T, V]]
}

func (s graphMUS[T, V]) Marshal(g Graph[T, V], bs []byte) (n int) {
	return s.vertecesMUS.Marshal(g.Vertices, bs)
}

func (s graphMUS[T, V]) Unmarshal(bs []byte) (g Graph[T, V], n int, err error) {
	g.Vertices, n, err = s.vertecesMUS.Unmarshal(bs)
	return
}

func (s graphMUS[T, V]) Size(g Graph[T, V]) (size int) {
	return s.vertecesMUS.Size(g.Vertices)
}

func (s graphMUS[T, V]) Skip(bs []byte) (n int, err error) {
	return s.vertecesMUS.Skip(bs)
}
