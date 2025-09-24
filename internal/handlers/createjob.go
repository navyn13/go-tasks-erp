package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func createjob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating Job")
	response := "Job created"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
