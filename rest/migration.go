package main

// MigrateProductV1 migrates V1 product to the current version.
func MigrateProductV1(productV1 ProductV1) (product Product, err error) {
	product = ProductV2{
		Name:        productV1.Name,
		Description: "Undefined",
	}
	return
}

// MigrateToProductV1 migrates current product to the V1 version.
func MigrateToProductV1(product Product) ProductV1 {
	return ProductV1{
		Name: product.Name,
	}
}

// MigrateAndMarshalProduct performs migration of the product to the specified
// version and marshals the result.
func MigrateAndMarshalProduct(dt DataType, product Product) (bs []byte) {
	switch dt {
	case ProductV1Type:
		productV1 := MigrateToProductV1(product)
		bs = make([]byte, SizeProductV1(productV1))
		MarshalProductV1(productV1, bs)
	case ProductV2Type:
		bs = make([]byte, SizeProductV2(product))
		MarshalProductV2(product, bs)
	}
	return
}

// UnmarshalAndMigrateProduct unmarshals from bs a product of the specified
// version and migrates the result to the current product version.
func UnmarshalAndMigrateProduct(dt DataType, bs []byte) (product Product,
	err error) {
	switch dt {
	case ProductV1Type:
		var productV1 ProductV1
		productV1, _, err = UnmarshalProductV1(bs)
		if err == nil {
			product, err = MigrateProductV1(productV1)
		}
	case ProductV2Type:
		product, _, err = UnmarshalProduct(bs)
	}
	return
}
