package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	response := "login got triggered"
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
