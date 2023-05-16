package main

import (
	"errors"

	muscom "github.com/mus-format/mus-common-go"
	"github.com/mus-format/mus-go/ord"
)

// Max length of the product name.
const NameMaxLength = 20

const (
	ProductV1Type DataType = iota + 1
	ProductV2Type
)

var (
	ProductType        = ProductV2Type
	MarshalMetaProduct = MarshalMetaProductV2
	SizeMetaProduct    = SizeMetaProductV2
	MarshalProduct     = MarshalProductV2
	UnmarshalProduct   = UnmarshalProductV2
	SizeProduct        = SizeProductV2
)

var ErrTooLongName = errors.New("too long name")

// -----------------------------------------------------------------------------
// Current product version.
type Product = ProductV2

// -----------------------------------------------------------------------------
type ProductV1 struct {
	Name string
}

func MarshalMetaProductV1(product ProductV1, bs []byte) (n int) {
	n = MarshalDataType(ProductV1Type, bs)
	return n + MarshalProductV1(product, bs[n:])
}

func SizeMetaProductV1(product ProductV1) (size int) {
	size += SizeDataType(ProductV1Type)
	return size + SizeProductV1(product)
}

func MarshalProductV1(product ProductV1, bs []byte) (n int) {
	return ord.MarshalString(product.Name, bs)
}

func UnmarshalProductV1(bs []byte) (product ProductV1, n int, err error) {
	var maxLength muscom.ValidatorFn[int] = func(length int) (err error) {
		if length > NameMaxLength {
			err = ErrTooLongName
		}
		return
	}
	product.Name, n, err = ord.UnmarshalValidString(maxLength, false, bs)
	return
}

func SizeProductV1(product ProductV1) (size int) {
	return ord.SizeString(product.Name)
}

// -----------------------------------------------------------------------------
type ProductV2 struct {
	Name        string
	Description string
}

func MarshalMetaProductV2(product ProductV2, bs []byte) (n int) {
	n = MarshalDataType(ProductV2Type, bs)
	return n + MarshalProductV2(product, bs[n:])
}

func SizeMetaProductV2(product ProductV2) (size int) {
	size += SizeDataType(ProductV2Type)
	return size + SizeProductV2(product)
}

func MarshalProductV2(product ProductV2, bs []byte) (n int) {
	n = ord.MarshalString(product.Name, bs)
	return n + ord.MarshalString(product.Description, bs[n:])
}

// UnmarshalProductV2 performs validation of the Product.Name field.
func UnmarshalProductV2(bs []byte) (product ProductV2, n int, err error) {
	var maxLength muscom.ValidatorFn[int] = func(length int) (err error) {
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

func SizeProductV2(product ProductV2) (size int) {
	size += ord.SizeString(product.Name)
	return size + ord.SizeString(product.Description)
}
