package main

type LinkedList[T any] struct {
	head *Elem[T]
	tail *Elem[T]
	len  int
}

func (l *LinkedList[T]) Head() *Elem[T] {
	return l.head
}

func (l *LinkedList[T]) Tail() *Elem[T] {
	return l.tail
}

func (l *LinkedList[T]) Len() int {
	return l.len
}

func (l *LinkedList[T]) AddFront(t T) {
	elem := Elem[T]{Val: t}
	if l.head == nil {
		l.tail = &elem
	} else {
		elem.next = l.head
		l.head.prev = &elem
	}
	l.head = &elem
	l.len++
}

func (l *LinkedList[T]) AddBack(t T) {
	elem := Elem[T]{Val: t}
	if l.tail == nil {
		l.head = &elem
	} else {
		elem.prev = l.tail
		l.tail.next = &elem
	}
	l.tail = &elem
	l.len++
}

func (l *LinkedList[T]) Remove(e *Elem[T]) {
	defer func() { l.len-- }()
	if e == l.head && e == l.tail {
		l.head = nil
		l.tail = nil
		return
	}
	if e == l.head {
		l.head = e.next
		e.next.prev = nil
		return
	}
	if e == l.tail {
		l.tail = e.prev
		e.prev.next = nil
		return
	}
	e.next.prev = e.prev
	e.prev.next = e.next
}
