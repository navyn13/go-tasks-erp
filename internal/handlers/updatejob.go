package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type updateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EmployeeID  int    `json:"employee_id"`
	JobID       int    `json:"job_id"`
}

func updateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Updating Job-------")

	adminUsernameCtx := r.Context().Value("username")
	adminRoleCtx := r.Context().Value("role")

	if adminUsernameCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := adminRoleCtx.(string)
	if !ok || role != "admin" {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	var req updateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.EmployeeID == 0 || req.JobID == 0 {
		http.Error(w, "Title, Description, EmployeeID, and JobID are required", http.StatusBadRequest)
		return
	}

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE jobs SET title = ?, description = ?, employee_id = ? WHERE id = ?",
		req.Title, req.Description, req.EmployeeID, req.JobID)
	if err != nil {
		http.Error(w, "Unable to update job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Job Updated Successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
