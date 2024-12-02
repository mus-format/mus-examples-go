package main

// Implements MarshallerMUS interface.
type Bar struct {
	str string
}

func (b Bar) MarshalMUS(bs []byte) (n int) {
	return MarshalBarMUS(b, bs) // Here BarDTS.Marshal() could be used.
}

func (b Bar) SizeMUS() (size int) {
	return SizeBarMUS(b)
}
