package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func startServer(products Products) {
	var (
		h = NewHandler(products)
		r = mux.NewRouter()
	)
	r.HandleFunc("/products/{id}", h.HandleGet).Methods(http.MethodGet)
	r.HandleFunc("/products/{id}", h.HandlePut).Methods(http.MethodPut)

	http.ListenAndServe(":8090", r)
}
