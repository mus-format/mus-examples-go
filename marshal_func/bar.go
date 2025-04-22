package main

// Bar implements the ext.MarshallerMUS interface.
type Bar struct {
	str string
}

func (b Bar) MarshalMUS(bs []byte) (n int) {
	return BarMUS.Marshal(b, bs)
}

func (b Bar) SizeMUS() (size int) {
	return BarMUS.Size(b)
}
