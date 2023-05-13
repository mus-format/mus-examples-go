package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/ymz-ncnk/assert"
)

func NewClient() Client {
	return Client{&http.Client{}}
}

// Client works with products of the current version.
type Client struct {
	client *http.Client
}

func (c Client) CreateProduct(id uuid.UUID, product Product) (err error) {
	bs := make([]byte, SizeProduct(product))
	MarshalProduct(product, bs)
	return createProduct(id, ProductType, bs, c.client)
}

func (c Client) GetProduct(id uuid.UUID) (product Product, err error) {
	bs, err := getProduct(id, ProductType, c.client)
	if err != nil {
		return
	}
	product, _, err = UnmarshalProduct(bs)
	return
}

// -----------------------------------------------------------------------------
func NewClientV1() ClientV1 {
	return ClientV1{&http.Client{}}
}

// ClientV1 works with products of the V1 version.
type ClientV1 struct {
	client *http.Client
}

func (c ClientV1) CreateProduct(id uuid.UUID, product ProductV1) (err error) {
	bs := make([]byte, SizeProductV1(product))
	MarshalProductV1(product, bs)
	return createProduct(id, ProductV1Type, bs, c.client)
}

func (c ClientV1) GetProduct(id uuid.UUID) (product ProductV1, err error) {
	bs, err := getProduct(id, ProductV1Type, c.client)
	if err != nil {
		return
	}
	product, _, err = UnmarshalProductV1(bs)
	return
}

// -----------------------------------------------------------------------------
// createProduct adds a "Data-Type" Header to the request, which allows the
// server to understand what version of the product is in the request.
func createProduct(id uuid.UUID, dt DataType, bs []byte, client *http.Client) (
	err error) {
	url := "http://localhost:8090/products/" + id.String()
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bs))
	if err != nil {
		return
	}

	req.Header.Add(DataTypeHeaderName, strconv.Itoa(int(dt)))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode == http.StatusBadRequest {
		defer resp.Body.Close()
		bs, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		return errors.New(string(bs))
	}
	assert.Equal(resp.StatusCode, http.StatusOK)
	return
}

// getProduct adds a "Data-Type" Header to the request, which helps the server
// understand which version of the product to return.
func getProduct(id uuid.UUID, dt DataType, client *http.Client) (bs []byte,
	err error) {
	url := "http://localhost:8090/products/" + id.String()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add(DataTypeHeaderName, strconv.Itoa(int(dt)))
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	bs, err = io.ReadAll(res.Body)
	return
}
