package main

import (
	"fmt"
	"io"

	"github.com/mus-format/mus-go/unsafe"
)

// In this example, we read and process several strings.
func main() {
	var (
		bs   = make([]byte, 10) // We make a long bs to fit all the read data.
		strs []string           // In this slice, we will accumulate all read strings.
		err  error
	)

	reader := MakeReader()
	for {
		_, err = reader.Read(bs)
		if err == io.EOF {
			break
		}
		// unsafe.UnmarshalString() creates a string that points to the given bs.
		// This means if we change bs after unmarshal, the content of the received
		// string will also change.

		// Here we use the same bs in each iteration. So we will receive strings
		// that point to the same bs.
		str, _, _ := unsafe.UnmarshalString(bs)

		// But this is not a problem if we process the received data before the next
		// read (which will change bs).

		fmt.Println(str) // Here, instead of fmt.Println(), we can, for example,
		// save the data to disk or send it over the network. In this case
		// everything is ok and the output will be:
		//
		// first
		// second
		//

		strs = append(strs, str) // But if we want to accumulate received strings,
		// we should use ord package instead, because in this case ...
	}

	fmt.Println(strs) // the output will be:
	//
	// [secon second]
	//

	// That's not strange, remember, with unsafe package all received strings
	// will point to the same bs (in this example), which at the end of the for
	// loop will equal to "second". The first string initially had a value of
	// "first" and a length of 5, so we see "secon" in this ouput.

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
		n = unsafe.MarshalString("first", p)
		return n, nil
	}).RegisterRead(func(p []byte) (n int, err error) {
		n = unsafe.MarshalString("second", p)
		return n, nil
	}).RegisterRead(func(p []byte) (n int, err error) {
		return 0, io.EOF
	})
}
