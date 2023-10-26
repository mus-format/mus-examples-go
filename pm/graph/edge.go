package main

type Edge[T any] struct {
	Weight int
	Vertex *Vertex[T]
}
