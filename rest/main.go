package main

import (
	"github.com/google/uuid"
	"github.com/ymz-ncnk/assert"
)

func init() {
	assert.On = true
}

// This example simulates a REST service that provides access to products.
// A little background. Initially, our products looked like this:
//
//	type Product struct {
//	  Name string
//	}
//
// But over time, it became necessary to add the "Description" field:
//
//	type Product struct {
//	  Name string
//	  Description string
//	}
//
// , and limit the length of the "Name" field to 20 characters.
// As a result, we have two versions of the products: V1 - old and V2 - current,
// as well as two clients: old and current.
func main() {
	var (
		oldID = uuid.New()
		// Products simulates the Persistence layer. Note that it only uses the
		// current version of the product.
		products = makeProducts(oldID)

		// In general, only the current product version should be used everywhere on
		// the server.
		//
		// To achieve this goal, we need to understand that there are only two
		// sources of old versions on the server - old clients and the storage.
		// And that, in the appropriate places (View and Persistence layers) we will
		// have to do the migration of old versions to the current one.
		// By the way, this migration looks like one fairly simple function.
		//
		// And one more thing, we also need to return old versions to old clients.
		// That is, we will still need to implement the migration of the current
		// version to the old one at the View level.
	)
	go startServer(products)

	var (
		client    = NewClient()
		oldClient = NewClientV1()
	)

	// The current client requests an older version of the product. The old
	// version on the server will be migrated to the current version, which will
	// be returned to the client.
	clientRequestsOldVersion(client, oldID)

	// An old client creates an old version of the product. The old version on the
	// server will be migrated to the current version, which will be saved to the
	// storage.
	oldClientCreatesOldVersion(oldClient)

	// The current client creates the current version of the product. The current
	// version is saved to the storage as is.
	id := clientCreatesCurrentVersion(client)

	// The old client requests the current version of the product. The current
	// version will be migrated on the server to the old version, which will be
	// returned to the client.
	oldClientRequestsCurrentVersion(oldClient, id)

	// The current client is trying to create an invalid product. An error will
	// occur during deserialization.
	clientCreatesInvalid(client)

	// An old client tries to create an invalid product. An error occurs when
	// migrating the old version to the current one.
	oldClientCreatesInvalid(oldClient)

}

func clientRequestsOldVersion(client Client, oldID uuid.UUID) {
	product, err := client.GetProduct(oldID)
	assert.EqualError(err, nil)
	assert.EqualDeep(product,
		Product{Name: "old", Description: "Undefined"})
}

func oldClientCreatesOldVersion(client ClientV1) {
	id := uuid.New()
	prodcut := ProductV1{Name: "another old"}
	err := client.CreateProduct(id, prodcut)
	assert.EqualError(err, nil)
}

func clientCreatesCurrentVersion(client Client) (id uuid.UUID) {
	id = uuid.New()
	product := Product{Name: "current", Description: "awesome"}
	err := client.CreateProduct(id, product)
	assert.EqualError(err, nil)
	return
}

func oldClientRequestsCurrentVersion(client ClientV1, id uuid.UUID) {
	product, err := client.GetProduct(id)
	assert.EqualError(err, nil)
	assert.EqualDeep(product, ProductV1{Name: "current"})
}

func oldClientCreatesInvalid(client ClientV1) {
	id := uuid.New()
	product := ProductV1{Name: "very very very long name"}
	err := client.CreateProduct(id, product)
	assert.EqualError(err, ErrTooLongName)
}

func clientCreatesInvalid(client Client) {
	id := uuid.New()
	product := Product{Name: "very very very long name"}
	err := client.CreateProduct(id, product)
	assert.EqualError(err, ErrTooLongName)
}

// makeProducts creates Products, which provides access to the old version of
// the product.
func makeProducts(idV1 uuid.UUID) Products {
	product := ProductV1{Name: "old"}
	bs := make([]byte, SizeMetaProductV1(product))
	MarshalMetaProductV1(product, bs)

	return NewProducts(map[uuid.UUID][]byte{idV1: bs})
}
