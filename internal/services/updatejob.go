package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func UpdateJob(w http.ResponseWriter, r *http.Request) {
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

	var req jobsSchema.UpdateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.EmployeeID == 0 || req.JobID == 0 {
		http.Error(w, "Title, Description, EmployeeID, and JobID are required", http.StatusBadRequest)
		return
	}

	db, err := db.GetDB()
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

	response := jobsSchema.UpdateJobResponse{Message: "Job Updated Successfully"}
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
