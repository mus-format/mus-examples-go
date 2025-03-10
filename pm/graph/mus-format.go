package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// MakeGraphSer creates a new Graph serializer that can be directly used.
func MakeGraphSer[T comparable, V any](keySer mus.Serializer[T],
	valSer mus.Serializer[V],
) mus.Serializer[Graph[T, V]] {
	var (
		ptrMap    = com.NewPtrMap()
		revPtrMap = com.NewReversePtrMap()
	)
	return pm.Wrap[Graph[T, V]](ptrMap, revPtrMap, NewGraphSer[T, V](ptrMap,
		revPtrMap, keySer, valSer))
}

// NewGraphSer creates a new Graph serializer.
func NewGraphSer[T comparable, V any](ptrMap *com.PtrMap,
	revPtrMap *com.ReversePtrMap,
	keySer mus.Serializer[T],
	valSer mus.Serializer[V],
) graphSer[T, V] {
	var (
		edgesSer    mus.Serializer[map[T]*Edge[T, V]]
		vertecesSer mus.Serializer[map[T]*Vertex[T, V]]

		edgePtrSer   mus.Serializer[*Edge[T, V]]
		vertexPtrSer mus.Serializer[*Vertex[T, V]]

		edgeSer   = edgeSer[T, V]{&vertexPtrSer}
		vertexSer = vertexSer[T, V]{valSer, &edgesSer}
	)
	edgePtrSer = pm.NewPtrSer[Edge[T, V]](ptrMap, revPtrMap, edgeSer)
	vertexPtrSer = pm.NewPtrSer[Vertex[T, V]](ptrMap, revPtrMap, vertexSer)

	edgesSer = ord.NewMapSer[T, *Edge[T, V]](keySer, edgePtrSer)
	vertecesSer = ord.NewMapSer[T, *Vertex[T, V]](keySer, vertexPtrSer)

	return graphSer[T, V]{vertecesSer}
}

// edgeSer implements the mus.Serializer interface for Edge.
type edgeSer[T comparable, V any] struct {
	vertexPtrSer *mus.Serializer[*Vertex[T, V]]
}

func (s edgeSer[T, V]) Marshal(e Edge[T, V], bs []byte) (n int) {
	n = varint.Int.Marshal(e.Weight, bs)
	return n + (*s.vertexPtrSer).Marshal(e.Vertex, bs[n:])
}

func (s edgeSer[T, V]) Unmarshal(bs []byte) (e Edge[T, V], n int, err error) {
	e.Weight, n, err = varint.Int.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	e.Vertex, n1, err = (*s.vertexPtrSer).Unmarshal(bs[n:])
	n += n1
	return
}

func (s edgeSer[T, V]) Size(e Edge[T, V]) (size int) {
	size = varint.Int.Size(e.Weight)
	return size + (*s.vertexPtrSer).Size(e.Vertex)
}

func (s edgeSer[T, V]) Skip(bs []byte) (n int, err error) {
	n, err = varint.Int.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.vertexPtrSer).Skip(bs[n:])
	n += n1
	return
}

// vertexSer implements the mus.Serializer interface for Vertex.
type vertexSer[T comparable, V any] struct {
	valSer   mus.Serializer[V]
	edgesSer *mus.Serializer[map[T]*Edge[T, V]]
}

func (s vertexSer[T, V]) Marshal(v Vertex[T, V], bs []byte) (n int) {
	n = s.valSer.Marshal(v.Val, bs)
	return n + (*s.edgesSer).Marshal(v.Edges, bs[n:])
}

func (s vertexSer[T, V]) Unmarshal(bs []byte) (v Vertex[T, V], n int, err error) {
	v.Val, n, err = s.valSer.Unmarshal(bs)
	if err != nil {
		return
	}
	var n1 int
	v.Edges, n1, err = (*s.edgesSer).Unmarshal(bs[n:])
	n += n1
	return
}

func (s vertexSer[T, V]) Size(v Vertex[T, V]) (size int) {
	size = s.valSer.Size(v.Val)
	return size + (*s.edgesSer).Size(v.Edges)
}

func (s vertexSer[T, V]) Skip(bs []byte) (n int, err error) {
	n, err = s.valSer.Skip(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = (*s.edgesSer).Skip(bs[n:])
	n += n1
	return
}

// graphSer implements the mus.Serializer interface for Graph.
type graphSer[T comparable, V any] struct {
	vertecesSer mus.Serializer[map[T]*Vertex[T, V]]
}

func (s graphSer[T, V]) Marshal(g Graph[T, V], bs []byte) (n int) {
	return s.vertecesSer.Marshal(g.Vertices, bs)
}

func (s graphSer[T, V]) Unmarshal(bs []byte) (g Graph[T, V], n int, err error) {
	g.Vertices, n, err = s.vertecesSer.Unmarshal(bs)
	return
}

func (s graphSer[T, V]) Size(g Graph[T, V]) (size int) {
	return s.vertecesSer.Size(g.Vertices)
}

func (s graphSer[T, V]) Skip(bs []byte) (n int, err error) {
	return s.vertecesSer.Skip(bs)
}
