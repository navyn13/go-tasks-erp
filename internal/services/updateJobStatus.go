package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func UpdateJobStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Updating Job Status--------")

	// FIX 1: Safely get user ID from context to prevent panic.
	userID, ok := r.Context().Value("id").(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusUnauthorized)
		return
	}

	// Decode request body first
	var req jobsSchema.UpdateJobStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var process string

	switch req.Process {
	case "cutting", "welding", "quality_check", "packaging", "dispatch":
		process = req.Process
	default:
		// If req.Process is not one of the allowed values, it's an invalid request.
		http.Error(w, "Invalid process name in request body", http.StatusBadRequest)
		return
	}

	switch req.Status {
	case "pending", "in-progress", "completed", "failed":

	default:
		// If req.Process is not one of the allowed values, it's an invalid request.
		http.Error(w, "Invalid Status given in request body", http.StatusBadRequest)
		return
	}

	// Also validate other fields after decoding
	if req.JobID == 0 {
		http.Error(w, "Invalid request: job_id is required", http.StatusBadRequest)
		return
	}

	// Get DB connection
	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Build the query safely using the validated column name
	updateProcessQuery := fmt.Sprintf(`
		UPDATE jobStatus js
		JOIN jobs j ON j.id = js.job_id
		SET js.%s = ?
		WHERE js.job_id = ? AND j.employee_id = ?`, process)

	res, err := dbConn.Exec(updateProcessQuery, req.Status, req.JobID, userID)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any row was updated
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "NO CHANGES MADE: either job is not assigned to you or process has that status before", http.StatusForbidden)
		return
	}

	// Response
	response := jobsSchema.UpdateJobStatusResponse{Message: "Job status updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
