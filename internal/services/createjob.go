package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Creating Job-------")

	// Get context and extract admin ID
	adminCtx := r.Context()
	adminID := adminCtx.Value("id")
	fmt.Println(adminID)

	var req jobsSchema.CreateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.EmployeeID == 0 {

		http.Error(w, "Title, Description and EmployeeID is required", http.StatusBadRequest)
		return
	}
	db, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec(
		"INSERT INTO jobs (title, description, employee_id, created_by_id) VALUES (?, ?, ?, ?)",
		req.Title, req.Description, req.EmployeeID, adminID,
	)
	if err != nil {
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jobID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get job ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create initial jobStatus row
	_, err = db.Exec("INSERT INTO jobStatus (job_id) VALUES (?)", jobID)
	if err != nil {
		http.Error(w, "Failed to create job status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := jobsSchema.CreateJobResponse{Message: "Job created successfully"}
	fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
