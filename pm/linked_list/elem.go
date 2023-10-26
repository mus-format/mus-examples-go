package main

type Elem[T any] struct {
	Val  T
	prev *Elem[T]
	next *Elem[T]
}

func (e Elem[T]) Next() *Elem[T] {
	return e.next
}

func (e Elem[T]) Prev() *Elem[T] {
	return e.prev
}
