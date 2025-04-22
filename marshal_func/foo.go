package main

// Foo implements the ext.MarshallerMUS interface.
type Foo struct {
	num int
}

func (f Foo) MarshalMUS(bs []byte) (n int) {
	return FooMUS.Marshal(f, bs)
}

func (f Foo) SizeMUS() (size int) {
	return FooMUS.Size(f)
}
