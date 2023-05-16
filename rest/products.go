package main

import (
	"errors"

	"github.com/google/uuid"
)

func NewProducts(m map[uuid.UUID][]byte) Products {
	return Products{m}
}

// Products imitates the Persistence layer. Take a note, it uses only the
// current version of the Product.
type Products struct {
	m map[uuid.UUID][]byte
}

// It saves the metadata (current Version) along with the product.
func (products Products) Add(id uuid.UUID, product Product) {
	bs := make([]byte, SizeMetaProduct(product))
	MarshalMetaProduct(product, bs)

	products.m[id] = bs
}

// Performs migration from the old Product version to the current one.
func (products Products) Get(id uuid.UUID) (product Product, err error) {
	bs, pst := products.m[id]
	if !pst {
		err = errors.New("not found")
		return
	}

	dt, n, err := UnmarshalDataType(bs)
	if err != nil {
		return
	}
	return UnmarshalAndMigrateProduct(dt, bs[n:])
}
