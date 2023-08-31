package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
// As a result, we have two versions of the products: ProductV1 - old and
// ProductV2 - current, as well as two clients - old and current.
//
// Only the current product version should be used everywhere on the server. To
// achieve this goal, we need to understand that there are only two sources of
// old versions on the server - old clients and the storage. And that, in the
// appropriate places (server handlers and persistence layer) we will have to do
// the migration of old versions to the current one with help of mus-dts-go and
// mus-dvs-go modules.
//
// You can see how this approach is being implemented in the following files:
// - client.go 		 	- client side.
// - server.go 		 	- contains server side handlers.
// - products.go 	 	- server side persistence layer.
// - product.go 	 	- defines product versions.
// - mus-format.go	- contains MUS format definitions.
func main() {
	var (
		oldID    = uuid.New()
		products = makeProducts(oldID)
	)
	go startServer(products)
	var (
		client    = NewClient()
		oldClient = NewOldClient()
	)
	// The current client requests an old version of the product. The old version
	// on the server will be migrated to the current version, which will be
	// returned to the client.
	clientRequestsOldVersion(client, oldID)

	// The old client creates an old version of the product. The old version on
	// the server will be migrated to the current version, which will be saved to
	// the storage.
	oldClientCreatesOldVersion(oldClient)

	// The current client creates the current version of the product. The current
	// version is saved to the storage as is.
	id := clientCreatesCurrentVersion(client)

	// The old client requests the current version of the product. The current
	// version will be migrated on the server to the old version, which will be
	// returned to the client.
	oldClientRequestsCurrentVersion(oldClient, id)

	// The current client tries to create an invalid product. An error will occur
	// on server during deserialization.
	clientCreatesInvalid(client)

	// The old client tries to create an invalid product. An error will occur
	// on server during deserialization.
	oldClientCreatesInvalid(oldClient)

}

func clientRequestsOldVersion(client Client, oldID uuid.UUID) {
	product, err := client.GetProduct(oldID)
	assert.EqualError(err, nil)
	assert.EqualDeep(product,
		Product{Name: "old", Description: "Undefined"})
}

func oldClientCreatesOldVersion(client OldClient) {
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

func oldClientRequestsCurrentVersion(client OldClient, id uuid.UUID) {
	product, err := client.GetProduct(id)
	assert.EqualError(err, nil)
	assert.EqualDeep(product, ProductV1{Name: "current"})
}

func oldClientCreatesInvalid(client OldClient) {
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
func makeProducts(id uuid.UUID) Products {
	product := ProductV1{Name: "old"}
	bs := make([]byte, ProductV1DTS.SizeMUS(product))
	ProductV1DTS.MarshalMUS(product, bs)
	return NewProducts(map[uuid.UUID][]byte{id: bs})
}

func startServer(products Products) {
	var (
		h = NewHandler(products)
		r = mux.NewRouter()
	)
	r.HandleFunc("/products/{id}", h.HandleGet).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.HandlePut).Methods(http.MethodPut)

	http.ListenAndServe(":8090", r)
}
