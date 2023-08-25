package main

// We have two versions of the product ProductV1 and ProductV2. The Product
// alias always indicates the current version of the product.

// Current product version.
type Product = ProductV2

type ProductV1 struct {
	Name string
}

type ProductV2 struct {
	Name        string
	Description string
}
