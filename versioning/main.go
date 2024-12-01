package main

import (
	"github.com/ymz-ncnk/assert"
)

// Demonstrates data versioning.
func main() {
	// Here we have two versions of Foo: FooV1 and FooV2 (current).
	fooV1 := FooV1{num: 10}

	// Let's marshal old V1 version using DTS.
	bs := make([]byte, FooV1DTS.Size(fooV1))
	FooV1DTS.Marshal(fooV1, bs)

	// Such bs, can be received from a legacy client, for example.

	// Now we can unmarshal it to the currect version.
	foo, _, err := UnmarshalFooMUS(bs) // UnmarshalFooMUS will migrate an
	// old version to the current.
	assert.EqualError(err, nil)
	assert.Equal(foo, Foo{str: "10"})

	// The current version can also be migrated to an older one if we want to send
	// it back to a legacy client.
}
