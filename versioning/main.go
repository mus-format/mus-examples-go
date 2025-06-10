package main

import (
	assert "github.com/ymz-ncnk/assert/panic"
)

func main() {
	// Marshal old V1 version using DTS.
	fooV1 := FooV1{num: 10}
	bs := make([]byte, FooV1DTS.Size(fooV1)) // Such bs, can be received from a
	// legacy client, for example.
	FooV1DTS.Marshal(fooV1, bs)

	// Unmarshal the current version from the bs.
	foo, _, err := FooMUS.Unmarshal(bs) // FooMUS.Unmarshal will migrate an
	// old version to the current.
	assert.EqualError(err, nil)
	assert.Equal(foo, Foo{str: "10"})

	// The current version can also be migrated to an older one if we want to send
	// it back to a legacy client.
}
