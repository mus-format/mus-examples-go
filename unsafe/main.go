package main

import (
	"fmt"
	"io"

	"github.com/mus-format/mus-go/unsafe"
)

// In this example, several strings are unmarshalled with the unsafe package.
func main() {
	var (
		bs   = make([]byte, 10) // Long enough bs to store all read data.
		strs []string           // This slice will accumulate strings.
		err  error
	)

	reader := MakeReader()
	// Read and process several strings.
	for {
		_, err = reader.Read(bs)
		if err == io.EOF {
			break
		}
		// unsafe.UnmarshalString() creates a string that points to the given bs.
		// This means that if bs changes later, the string's content will also
		// change.

		// The same bs is used in each iteration, so all received strings are
		// point to it.
		str, _, _ := unsafe.UnmarshalString(nil, bs)

		// This is not a problem if each string is processed before the next read
		// (which changes bs).
		fmt.Println(str) // Instead of using fmt.Println(), str, for example, can be
		// saved to the disk or sent over the network.
		//
		// The output will be:
		//
		// first
		// second
		//

		// But if we want to accumulate received strings, the ord package should be
		// used instead, because in this case ...
		strs = append(strs, str)
	}

	fmt.Println(strs)
	// ... the output will be:
	//
	// [secon second]
	//

	// That's not strange, with unsafe package all received strings will point to
	// the same bs, which at the end of the for loop equals to "second". The first
	// string initially had a value of "first" and a length of 5, so we see
	// "secon" in the ouput.

	// The unsafe package provides high performance but requires careful use when
	// it comes to strings. The good news - with other types there is no such
	// behavior.

	// And one more thing. If each string will have own bs (we can create new one
	// at each iteration), everything will be ok and the output of the
	// fmt.Println(strs) will be:
	//
	// [first second]
	//
}

// MakeReader creates an io.Reader that returns “first” on the first read and
// “second” on the second. All results are marshalled with an unsafe package.
func MakeReader() io.Reader {
	return NewReaderMock().RegisterRead(func(p []byte) (n int, err error) {
		n = unsafe.MarshalString("first", nil, p)
		return n, nil
	}).RegisterRead(func(p []byte) (n int, err error) {
		n = unsafe.MarshalString("second", nil, p)
		return n, nil
	}).RegisterRead(func(p []byte) (n int, err error) {
		return 0, io.EOF
	})
}
