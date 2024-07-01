package main

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/pm"
	"github.com/mus-format/mus-go/varint"
)

// -----------------------------------------------------------------------------
// Edge
// -----------------------------------------------------------------------------

// NewEdgeMarshaller creates a new Edge marshaller.
//
// vertexM param is a Vertex marshaller.
func NewEdgeMarshaller[T any](
	vertexM *mus.Marshaller[*Vertex[T]]) mus.Marshaller[Edge[T]] {
	return mus.MarshallerFn[Edge[T]](
		func(e Edge[T], bs []byte) (n int) {
			n = varint.MarshalInt(e.Weight, bs)
			return n + (*vertexM).MarshalMUS(e.Vertex, bs[n:])
		},
	)
}

// NewEdgeUnmarshaller cretes a new Edge unmarshaller.
//
// vertexU param is a Vertex unmarshaller.
func NewEdgeUnmarshaller[T any](
	vertexU *mus.Unmarshaller[*Vertex[T]]) mus.Unmarshaller[Edge[T]] {
	return mus.UnmarshallerFn[Edge[T]](
		func(bs []byte) (e Edge[T], n int, err error) {
			e.Weight, n, err = varint.UnmarshalInt(bs)
			if err != nil {
				return
			}
			var n1 int
			e.Vertex, n1, err = (*vertexU).UnmarshalMUS(bs[n:])
			n += n1
			return
		},
	)
}

// NewEdgeSizer creates a new Edge sizer.
//
// vertexS param is a Vertex sizer.
func NewEdgeSizer[T any](vertexS *mus.Sizer[*Vertex[T]]) mus.Sizer[Edge[T]] {
	return mus.SizerFn[Edge[T]](
		func(e Edge[T]) (size int) {
			size = varint.SizeInt(e.Weight)
			return size + (*vertexS).SizeMUS(e.Vertex)
		},
	)
}

// -----------------------------------------------------------------------------
// Vertex
// -----------------------------------------------------------------------------

// NewVertexMarshaller creates a new Vertex marshaller.
//
// valM param is a Vertex value marshaller, edgeM - Edge marshaller.
func NewVertexMarshaller[T any](valM mus.Marshaller[T],
	edgeM *mus.Marshaller[*Edge[T]]) mus.Marshaller[Vertex[T]] {
	return mus.MarshallerFn[Vertex[T]](
		func(v Vertex[T], bs []byte) (n int) {
			n = valM.MarshalMUS(v.Val, bs)
			return n + ord.MarshalMap[int, *Edge[T]](v.Edges, nil,
				mus.MarshallerFn[int](varint.MarshalInt),
				*edgeM,
				bs[n:])
		},
	)
}

// NewVertexUnmarshaller creates a new Vertex unmarshaller.
//
// valU param is a Vertex value unmarshaller, edgeU - Edge unmarshaller.
func NewVertexUnmarshaller[T any](valU mus.Unmarshaller[T],
	edgeU *mus.Unmarshaller[*Edge[T]]) mus.Unmarshaller[Vertex[T]] {
	return mus.UnmarshallerFn[Vertex[T]](
		func(bs []byte) (t Vertex[T], n int, err error) {
			t.Val, n, err = valU.UnmarshalMUS(bs)
			// t.Val, n, err = varint.UnmarshalInt(bs)
			if err != nil {
				return
			}
			var n1 int
			t.Edges, n1, err = ord.UnmarshalMap[int, *Edge[T]](nil,
				mus.UnmarshallerFn[int](varint.UnmarshalInt), *edgeU, bs[n:])
			n += n1
			return
		},
	)
}

// NewVertexSizer creates a new Vertex sizer.
//
// valS param is a Vertex value sizer, edgeS - Edge sizer.
func NewVertexSizer[T any](valS mus.Sizer[T],
	edgeS *mus.Sizer[*Edge[T]]) mus.Sizer[Vertex[T]] {
	return mus.SizerFn[Vertex[T]](
		func(v Vertex[T]) (size int) {
			size = valS.SizeMUS(v.Val)
			return size + ord.SizeMap[int, *Edge[T]](v.Edges, nil,
				mus.SizerFn[int](varint.SizeInt), *edgeS)
		},
	)
}

// -----------------------------------------------------------------------------
// Graph
// -----------------------------------------------------------------------------

// NewGraphMarshaller creates a new Graph marshaller.
//
// keyM param is a Vertex key marshaller, valM - Vertex value marshaller.
func NewGraphMarshaller[T comparable, V any](keyM mus.Marshaller[T],
	valM mus.Marshaller[V]) mus.Marshaller[Graph[T, V]] {
	var (
		mp                  = com.NewPtrMap()
		vertexPtrMarshaller mus.Marshaller[*Vertex[V]]
		edgePtrMarshaller   mus.Marshaller[*Edge[V]]
	)
	vertexPtrMarshaller = mus.MarshallerFn[*Vertex[V]](
		func(v *Vertex[V], bs []byte) (n int) {
			return pm.MarshalPtr[Vertex[V]](v, NewVertexMarshaller(valM,
				&edgePtrMarshaller),
				mp,
				bs)
		},
	)
	edgePtrMarshaller = mus.MarshallerFn[*Edge[V]](
		func(e *Edge[V], bs []byte) (n int) {
			return ord.MarshalPtr[Edge[V]](e, NewEdgeMarshaller(&vertexPtrMarshaller),
				bs)
		},
	)
	return mus.MarshallerFn[Graph[T, V]](
		func(g Graph[T, V], bs []byte) (n int) {
			return ord.MarshalMap[T, *Vertex[V]](g.Vertices, nil, keyM,
				vertexPtrMarshaller,
				bs)
		},
	)
}

// NewGraphUnmarshaller creates a new Graph unmarshaller.
//
// keyU param is a Vertex key unmarshaller, valU - Vertex value unmarshaller.
func NewGraphUnmarshaller[T comparable, V any](keyU mus.Unmarshaller[T],
	valU mus.Unmarshaller[V]) mus.Unmarshaller[Graph[T, V]] {
	var (
		mp                    = com.NewReversePtrMap()
		vertexPtrUnmarshaller mus.Unmarshaller[*Vertex[V]]
		edgePtrUnmarshaller   mus.Unmarshaller[*Edge[V]]
	)
	vertexPtrUnmarshaller = mus.UnmarshallerFn[*Vertex[V]](
		func(bs []byte) (v *Vertex[V], n int, err error) {
			return pm.UnmarshalPtr[Vertex[V]](NewVertexUnmarshaller(valU,
				&edgePtrUnmarshaller),
				mp,
				bs)
		},
	)
	edgePtrUnmarshaller = mus.UnmarshallerFn[*Edge[V]](
		func(bs []byte) (e *Edge[V], n int, err error) {
			return ord.UnmarshalPtr[Edge[V]](
				NewEdgeUnmarshaller(&vertexPtrUnmarshaller),
				bs)
		},
	)
	return mus.UnmarshallerFn[Graph[T, V]](
		func(bs []byte) (g Graph[T, V], n int, err error) {
			g.Vertices, n, err = ord.UnmarshalMap[T, *Vertex[V]](nil, keyU,
				vertexPtrUnmarshaller,
				bs)
			return
		},
	)
}

// NewGraphSizer creates a new Graph sizer.
//
// keyU param is a Vertex key sizer, valU - Vertex value sizer.
func NewGraphSizer[T comparable, V any](keyS mus.Sizer[T],
	valS mus.Sizer[V]) mus.Sizer[Graph[T, V]] {
	var (
		mp             = com.NewPtrMap()
		vertexPtrSizer mus.Sizer[*Vertex[V]]
		edgePtrSizer   mus.Sizer[*Edge[V]]
	)
	vertexPtrSizer = mus.SizerFn[*Vertex[V]](
		func(v *Vertex[V]) (size int) {
			return pm.SizePtr[Vertex[V]](v, NewVertexSizer(valS, &edgePtrSizer), mp)
		},
	)
	edgePtrSizer = mus.SizerFn[*Edge[V]](
		func(e *Edge[V]) (size int) {
			return ord.SizePtr[Edge[V]](e, NewEdgeSizer(&vertexPtrSizer))
		},
	)
	return mus.SizerFn[Graph[T, V]](
		func(g Graph[T, V]) (size int) {
			return ord.SizeMap[T, *Vertex[V]](g.Vertices, nil, keyS, vertexPtrSizer)
		},
	)
}
