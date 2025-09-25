package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func DeleteJob(w http.ResponseWriter, r *http.Request) {
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

	var req jobsSchema.DeleteJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.GetDB()
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

	response := jobsSchema.DeleteJobResponse{Message: "Job deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
