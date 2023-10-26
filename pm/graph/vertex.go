package main

type Vertex[T any] struct {
	Val   T
	Edges map[int]*Edge[T]
}
