package main

import "github.com/ymz-ncnk/mok"

type ReadFn func(p []byte) (n int, err error)

func NewReaderMock() ReaderMock {
	return ReaderMock{mok.New("Reader")}
}

// ReaderMock is a mock implementation of the io.Reader. It simply uses
// mok.Mock as a delegate.
type ReaderMock struct {
	*mok.Mock
}

// RegisterRead registers a function as a single call to the Read() method.
func (m ReaderMock) RegisterRead(fn ReadFn) ReaderMock {
	m.Register("Read", fn)
	return m
}

func (m ReaderMock) Read(p []byte) (n int, err error) {
	result, err := m.Call("Read", p)
	if err != nil {
		return
	}
	n = result[0].(int)
	err, _ = result[1].(error)
	return
}
