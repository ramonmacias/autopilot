package api

import "net/http"

func ShowContact(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func CreateContact(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
