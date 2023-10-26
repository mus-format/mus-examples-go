package main

type Graph[T comparable, V any] struct {
	Vertices map[T]*Vertex[V]
}
