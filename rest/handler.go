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
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	assert.EqualError(err, nil)

	product, err := h.products.Get(id)
	assert.EqualError(err, nil)

	dt, err := strconv.ParseUint(r.Header.Get(DataTypeHeaderName), 10, 8)
	assert.EqualError(err, nil)

	bs := MigrateAndMarshalProduct(DataType(dt), product)

	w.Write(bs)
}

// HandlePut performs migration from the old Product version to the current
// one.
func (h Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	assert.EqualError(err, nil)

	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	assert.EqualError(err, nil)

	dt, err := strconv.ParseUint(r.Header.Get(DataTypeHeaderName), 10, 8)
	assert.EqualError(err, nil)

	product, err := UnmarshalAndMigrateProduct(DataType(dt), bs)
	if err == ErrTooLongName {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	assert.EqualError(err, nil)

	h.products.Add(id, product)
}
