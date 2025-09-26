package utils

import (
	"net/http"
)

// GetHeader retrieves the value of a specified header from an HTTP request.
func GetHeader(r *http.Request, headerName string) string {
	return r.Header.Get(headerName)
}
