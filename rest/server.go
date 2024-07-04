package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	com "github.com/mus-format/common-go"
	"github.com/ymz-ncnk/assert"
)

// On the server side, we use the mus-dvs-go module, namely ProductDVS, to:
// 1. Get the current version of the product from any client.
// 2. Send to the client the version of the product it needs.

// With this header, the server handler finds out which version of the product
// should be returned with the response to the GET request.
const DTMHeaderName = "DTM"

func NewHandler(products Products) Handler {
	return Handler{products}
}

// For the sake of simplicity, it directly uses Products.
type Handler struct {
	products Products
}

// HandleGet handles GET requests.
func (h Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Receives the current version of the product from the storage.
	id := parseID(r)
	product, err := h.products.Get(id)
	assert.EqualError(err, nil)

	// Using DTM, the client indicates which version of the product it wants to
	// receive.
	dtm := parseDTM(r)
	// ProducDVS will migrate the product to the appropriate version.
	bs, _, err := ProductDVS.MakeBSAndMarshal(dtm, product)
	assert.EqualError(err, nil)

	_, err = w.Write(bs)
	assert.EqualError(err, nil)
}

// HandlePut handles PUT requests.
func (h Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var (
		id = parseID(r)
		bs = readRequestBody(r)
	)
	// bs contains a DTM and data itself. The DTM also determines the version of
	// the data.
	//
	// If there is an old version of the product in bs, ProductDVS migrates it to
	// the current version.
	_, product, _, err := ProductDVS.Unmarshal(bs)
	if err != nil {
		sendBackUnmarshallerr(w, err)
	}
	// This way, only the current version of the product is always saved to the
	// storage.
	h.products.Add(id, product)
}

func parseID(r *http.Request) (id uuid.UUID) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	assert.EqualError(err, nil)
	return
}

func parseDTM(r *http.Request) (dtm com.DTM) {
	n, err := strconv.ParseUint(r.Header.Get(DTMHeaderName), 10, 8)
	assert.EqualError(err, nil)
	return com.DTM(n)
}

func readRequestBody(r *http.Request) []byte {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	assert.EqualError(err, nil)
	return bs
}

func sendBackUnmarshallerr(w http.ResponseWriter, err error) {
	if err == ErrTooLongName {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	assert.EqualError(err, nil)
}
