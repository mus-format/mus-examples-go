package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ymz-ncnk/assert"
)

// With this header, the handler finds out which version of the product is in
// the PUT request, or which version of the product should be returned with the
// response to the GET request.
const DataTypeHeaderName = "Data-Type"

func NewHandler(products Products) Handler {
	return Handler{products}
}

// For the sake of simplicity, it directly uses Products.
type Handler struct {
	products Products
}

// HandleGet performs migration from the current Product version to the old one.
func (h Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := parseID(r)
	product, err := h.products.Get(id)
	assert.EqualError(err, nil)

	var (
		dt = parseDataType(r)
		bs = MigrateAndMarshalProduct(dt, product)
	)
	_, err = w.Write(bs)
	assert.EqualError(err, nil)
}

// HandlePut performs migration from the old Product version to the current
// one.
func (h Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var (
		id = parseID(r)
		bs = readRequestBody(r)
		dt = parseDataType(r)
	)
	product, err := UnmarshalAndMigrateProduct(dt, bs)
	if err != nil {
		sendBackUnmarshalErr(w, err)
	}
	h.products.Add(id, product)
}

func parseID(r *http.Request) (id uuid.UUID) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	assert.EqualError(err, nil)
	return
}

func parseDataType(r *http.Request) (dt DataType) {
	n, err := strconv.ParseUint(r.Header.Get(DataTypeHeaderName), 10, 8)
	assert.EqualError(err, nil)
	return DataType(n)
}

func readRequestBody(r *http.Request) []byte {
	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	assert.EqualError(err, nil)
	return bs
}

func sendBackUnmarshalErr(w http.ResponseWriter, err error) {
	if err == ErrTooLongName {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	assert.EqualError(err, nil)
}
