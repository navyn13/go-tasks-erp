package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func signup(w http.ResponseWriter, r *http.Request) {
	response := "signup got triggered"
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
