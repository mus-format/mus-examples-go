package main

type Vertex[T comparable, V any] struct {
	Val   V
	Edges map[T]*Edge[T, V]
}
