package main

import (
	"fmt"

	"github.com/mus-format/mus-go/varint"
)

func main() {
	// Encode three numbers in turn - 5, 10, 15.
	bs := make([]byte, varint.Int.Size(5)+varint.Int.Size(10)+varint.Int.Size(15))
	n := varint.Int.Marshal(5, bs)
	n += varint.Int.Marshal(10, bs[n:])
	varint.Int.Marshal(15, bs[n:])

	// Get them back in the opposite direction. Errors are omitted for simplicity.
	n1, _ := varint.Int.Skip(bs)
	n2, _ := varint.Int.Skip(bs)
	num, _, _ := varint.Int.Unmarshal(bs[n1+n2:])
	fmt.Println(num)
	num, _, _ = varint.Int.Unmarshal(bs[n1:])
	fmt.Println(num)
	num, _, _ = varint.Int.Unmarshal(bs)
	fmt.Println(num)
	// The output will be:
	// 15
	// 10
	// 5
}
