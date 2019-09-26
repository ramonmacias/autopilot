package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func BuildRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/contact/{id}", ShowContact).Methods("GET")
	r.HandleFunc("/contact/{id}", UpdateContact).Methods("PUT")
	r.HandleFunc("/contact", CreateContact).Methods("POST")

	http.Handle("/", r)
	return r
}
