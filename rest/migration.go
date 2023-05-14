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
