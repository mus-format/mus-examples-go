package main

type FooV1 struct {
	num int
}

func (f FooV1) MarshalMUS(bs []byte) (n int) {
	return FooV1DTS.Marshal(f, bs)
}

func (f FooV1) SizeMUS() (size int) {
	return FooV1DTS.Size(f)
}
