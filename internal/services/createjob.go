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
	adminIDInterface := adminCtx.Value("id")
	adminID, ok := adminIDInterface.(int)
	if !ok {
		http.Error(w, "Invalid admin ID in context", http.StatusInternalServerError)
		return
	}

	var req jobsSchema.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.EmployeeID == 0 {
		http.Error(w, "Title, Description and EmployeeID are required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Check if assigned employee is an admin
	var adminCount int
	err = dbConn.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND role = 'admin'", req.EmployeeID).Scan(&adminCount)
	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if adminCount > 0 {
		http.Error(w, "Cannot assign a job to an admin user", http.StatusForbidden)
		return
	}

	// Start transaction
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, "Failed to start DB transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert job
	result, err := tx.Exec(
		"INSERT INTO jobs (title, description, employee_id, created_by_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
		req.Title, req.Description, req.EmployeeID, adminID,
	)
	if err != nil {
		tx.Rollback()
		http.Error(w, "DB insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jobID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to get job ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	processes := []string{"cutting", "welding", "qualitycheck", "packaging", "dispatch"}
	for _, p := range processes {
		_, err := tx.Exec("INSERT INTO jobStatus (jobid, process_name, status) VALUES (?, ?, 'pending')", jobID, p)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create job processes: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := jobsSchema.CreateJobResponse{Message: "Job created successfully"}
	fmt.Println(response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
