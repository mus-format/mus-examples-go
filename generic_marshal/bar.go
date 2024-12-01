package main

// Implements MarshallerMUS interface. Could implement it with DTS instead own
// Marshal/Size functions.
type Bar struct {
	str string
}

func (b Bar) MarshalMUS(bs []byte) (n int) {
	return MarshalBarMUS(b, bs) // Here BarDTS.Marshal() could be used.
}

func (b Bar) SizeMUS() (size int) {
	return SizeBarMUS(b)
}
