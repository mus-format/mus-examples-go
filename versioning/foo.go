package main

// Indicates the current version.
type Foo FooV2

func (f Foo) MarshalMUS(bs []byte) (n int) {
	return FooMUS.Marshal(f, bs)
}

func (f Foo) SizeMUS() (size int) {
	return FooMUS.Size(f)
}
