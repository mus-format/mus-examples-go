package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	com "github.com/mus-format/common-go"
	"github.com/ymz-ncnk/assert"
)

// Each client works with only one product version, so we use the mus-dts-go
// module (which helps us to marshal/unmarshal specific product version) on the
// client side.

// -----------------------------------------------------------------------------
// Current client.
// -----------------------------------------------------------------------------
func NewClient() Client {
	return Client{&http.Client{}}
}

// Client works with products of the current version.
type Client struct {
	client *http.Client
}

func (c Client) CreateProduct(id uuid.UUID, product Product) (err error) {
	// ProductDTS encodes product DTM (which also indicates the version of the
	// product) and product itself to the bs. This allows the server to know which
	// version of the product it is dealing with.
	bs := make([]byte, ProductDTS.SizeMUS(product))
	ProductDTS.MarshalMUS(product, bs)
	return doPUT(id, bs, c.client)
}

func (c Client) GetProduct(id uuid.UUID) (product Product, err error) {
	// We need to specify which version of the product we want to receive from the
	// server, so we send ProductDTM in the request header.
	bs, err := doGet(id, ProductDTM, c.client)
	if err != nil {
		return
	}
	// bs here will also contain the product DTM and the product itself, so we
	// should use ProductDTS.
	product, _, err = ProductDTS.UnmarshalMUS(bs)
	return
}

// -----------------------------------------------------------------------------
// Old client.
// -----------------------------------------------------------------------------
func NewOldClient() OldClient {
	return OldClient{&http.Client{}}
}

// OldClient works with products of the old V1 version. Uses ProductV1DTS.
type OldClient struct {
	client *http.Client
}

func (c OldClient) CreateProduct(id uuid.UUID, product ProductV1) (err error) {
	bs := make([]byte, ProductV1DTS.SizeMUS(product))
	ProductV1DTS.MarshalMUS(product, bs)
	return doPUT(id, bs, c.client)
}

func (c OldClient) GetProduct(id uuid.UUID) (product ProductV1, err error) {
	bs, err := doGet(id, ProductV1DTM, c.client)
	if err != nil {
		return
	}
	product, _, err = ProductV1DTS.UnmarshalMUS(bs)
	return
}

// -----------------------------------------------------------------------------
// Helper functions.
// -----------------------------------------------------------------------------
func doPUT(id uuid.UUID, bs []byte, client *http.Client) (
	err error) {
	url := "http://localhost:8090/products/" + id.String()
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bs))
	if err != nil {
		return
	}
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

// doGet adds a "DTM" header to the request, which helps the server understand
// what version of the product it should return.
func doGet(id uuid.UUID, dtm com.DTM, client *http.Client) (bs []byte,
	err error) {
	url := "http://localhost:8090/products/" + id.String()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add(DTMHeaderName, strconv.Itoa(int(dtm)))
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	bs, err = io.ReadAll(res.Body)
	return
}
