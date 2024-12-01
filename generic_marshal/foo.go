package main

// Implements MarshallerMUS interface. Could implement it with DTS instead own
// Marshal/Size functions.
type Foo struct {
	num int
}

func (f Foo) MarshalMUS(bs []byte) (n int) {
	return MarshalFooMUS(f, bs) // Here FooDTS.Marshal() could be used.
}

func (f Foo) SizeMUS() (size int) {
	return SizeFooMUS(f)
}
