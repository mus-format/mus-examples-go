package main

// Foo indicates the current version.
type Foo FooV2

type FooV1 struct {
	num int
}

type FooV2 struct {
	str string
}
