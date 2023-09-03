package main

import (
	"errors"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	dvs "github.com/mus-format/mus-dvs-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
)

// Contains all MUS format related constants, variables and functions. Here you
// can find:
// - Marshal/Unmarshal/Size functions for each product version.
// - Data Type Metadata Support (DTS) for each product version.
// - Data Versioning Support (DVS) for the product.

// -----------------------------------------------------------------------------
// DTM
// -----------------------------------------------------------------------------

// Each DTM also defines a product version.
const (
	ProductV1DTM com.DTM = iota
	ProductV2DTM
)

// DTM of the current product version.
const ProductDTM = ProductV2DTM

// -----------------------------------------------------------------------------
// Marshal/Unmarshal/Size functions
// -----------------------------------------------------------------------------

// Max length of the product name (used for validation).
const NameMaxLength = 20

var ErrTooLongName = errors.New("too long name")

// Marshal function for the first product version.
func MarshalProductV1MUS(product ProductV1, bs []byte) (n int) {
	return ord.MarshalString(product.Name, bs)
}

// Unmarshal function for the first product version.
//
// UnmarshalProductV1MUS performs validation of the Product.Name field.
func UnmarshalProductV1MUS(bs []byte) (product ProductV1, n int, err error) {
	var maxLength com.ValidatorFn[int] = func(length int) (err error) {
		if length > NameMaxLength {
			err = ErrTooLongName
		}
		return
	}
	product.Name, n, err = ord.UnmarshalValidString(maxLength, false, bs)
	return
}

// Size function for the first product version.
func SizeProductV1MUS(product ProductV1) (size int) {
	return ord.SizeString(product.Name)
}

// Marshal function for the second product version.
func MarshalProductV2MUS(product ProductV2, bs []byte) (n int) {
	n = ord.MarshalString(product.Name, bs)
	return n + ord.MarshalString(product.Description, bs[n:])
}

// Unmarshal function for the second product version.
//
// UnmarshalProductV2MUS performs validation of the Product.Name field.
func UnmarshalProductV2MUS(bs []byte) (product ProductV2, n int, err error) {
	var maxLength com.ValidatorFn[int] = func(length int) (err error) {
		if length > NameMaxLength {
			err = ErrTooLongName
		}
		return
	}
	product.Name, n, err = ord.UnmarshalValidString(maxLength, false, bs)
	if err != nil {
		return
	}
	product.Description, _, err = ord.UnmarshalString(bs[n:])
	return
}

// Size function for the second product version.
func SizeProductV2MUS(product ProductV2) (size int) {
	size += ord.SizeString(product.Name)
	return size + ord.SizeString(product.Description)
}

// -----------------------------------------------------------------------------
// Data Type Metadata Support (DTS)
// -----------------------------------------------------------------------------

var (
	// DTS of the first product version.
	ProductV1DTS = dts.New[ProductV1](ProductV1DTM,
		mus.MarshallerFn[ProductV1](MarshalProductV1MUS),
		mus.UnmarshallerFn[ProductV1](UnmarshalProductV1MUS),
		mus.SizerFn[ProductV1](SizeProductV1MUS))
	// DTS of the second product version.
	ProductV2DTS = dts.New[ProductV2](ProductV2DTM,
		mus.MarshallerFn[ProductV2](MarshalProductV2MUS),
		mus.UnmarshallerFn[ProductV2](UnmarshalProductV2MUS),
		mus.SizerFn[ProductV2](SizeProductV2MUS))
)

// DTS of the current product version.
var ProductDTS = ProductV2DTS

// -----------------------------------------------------------------------------
// Data Versioning Support (DVS)
// -----------------------------------------------------------------------------

var registry = com.NewRegistry([]com.TypeVersion{
	dvs.Version[ProductV1, Product]{
		DTS: ProductV1DTS,
		MigrateOld: func(t ProductV1) (v Product, err error) {
			v = ProductV2{
				Name:        t.Name,
				Description: "Undefined",
			}
			return
		},
		MigrateCurrent: func(v Product) (t ProductV1, err error) {
			t = ProductV1{
				Name: v.Name,
			}
			return
		},
	},
	dvs.Version[ProductV2, Product]{
		DTS: ProductV2DTS,
		MigrateOld: func(t ProductV2) (v ProductV2, err error) {
			return t, nil
		},
		MigrateCurrent: func(v ProductV2) (t ProductV2, err error) {
			return v, nil
		},
	},
})

var ProductDVS = dvs.New[Product](registry)
