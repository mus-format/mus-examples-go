package main

// Foo implements the MarshallerMUS interface.
type Foo struct {
	num int
}

func (f Foo) MarshalMUS(bs []byte) (n int) {
	return FooMUS.Marshal(f, bs) // Here FooDTS.Marshal() could be used.
}

func (f Foo) SizeMUS() (size int) {
	return FooMUS.Size(f)
}
