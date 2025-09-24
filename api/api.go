package api

import (
	"net/http"
)

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalServerErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusInternalServerError)
	}
)
