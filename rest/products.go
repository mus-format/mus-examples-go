package main

import (
	"errors"

	"github.com/google/uuid"
)

// Products simulates a persistence layer. Just as the server can receive old
// product versions from the old clients, Products can receive them from the
// storage. Therefore, we also have to use the mus-dvs-go module here.

func NewProducts(m map[uuid.UUID][]byte) Products {
	return Products{m}
}

// Take a note, Products accepts and returns only the current version of
// product.
type Products struct {
	m map[uuid.UUID][]byte
}

// When saving a product, we must also save its version, so we use ProductDTMS
// here.
func (products Products) Add(id uuid.UUID, product Product) {
	bs := make([]byte, ProductDTMS.SizeMUS(product))
	ProductDTMS.MarshalMUS(product, bs)
	products.m[id] = bs
}

// We can get an old version of a product from the storage, so we use ProductDVS
// here, which migrates such old versions to the current one.
func (products Products) Get(id uuid.UUID) (product Product, err error) {
	bs, pst := products.m[id]
	if !pst {
		err = errors.New("not found")
		return
	}
	_, product, _, err = ProductDVS.UnmarshalMUS(bs)
	return
}
