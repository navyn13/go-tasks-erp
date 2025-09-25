package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteJobRequest struct {
	JobID int `json:"job_id"`
}

func deletejob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Deleting Job-------")
	adminUsernameCtx := r.Context().Value("username")
	adminRoleCtx := r.Context().Value("role")
	if adminRoleCtx != "admin" {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}
	if adminUsernameCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req DeleteJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM jobs WHERE id = ?", req.JobID)
	if err != nil {
		http.Error(w, "DB delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := "Job deleted"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
