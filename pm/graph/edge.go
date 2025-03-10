package main

type Edge[T comparable, V any] struct {
	Weight int
	Vertex *Vertex[T, V]
}
