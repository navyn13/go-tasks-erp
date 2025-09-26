package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func UpdateJobStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Updating Job Status--------")

	userID, ok := r.Context().Value("id").(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusUnauthorized)
		return
	}

	var req jobsSchema.UpdateJobStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

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

	query := "UPDATE jobStatus js JOIN jobs j ON js.jobid = j.id SET js.status = ?"
	args := []interface{}{req.Status}

	if req.Status == "in-progress" {
		query += ", js.started_at = ?"
		args = append(args, time.Now())
	} else if req.Status == "completed" {
		query += ", js.completed_at = ?"
		args = append(args, time.Now())
	}

	query += " WHERE js.jobid = ? AND js.process_name = ? AND j.employee_id = ?"
	args = append(args, req.JobID, req.Process, userID)

	res, err := dbConn.Exec(query, args...)

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
