package main

// Bar implements the MarshallerMUS interface.
type Bar struct {
	str string
}

func (b Bar) MarshalMUS(bs []byte) (n int) {
	return BarMUS.Marshal(b, bs) // Here BarDTS.Marshal() could be used.
}

func (b Bar) SizeMUS() (size int) {
	return BarMUS.Size(b)
}
