package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	fmt.Println("--------Getting Job Status-------")

	var req jobsSchema.GetJobStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	db, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query("SELECT job_id, cutting, welding, quality_check, packaging, dispatch FROM jobStatus WHERE job_id = ?", req.JobID)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobStatuses []jobsSchema.GetJobStatusResponse
	for rows.Next() {
		var jobStatus jobsSchema.GetJobStatusResponse
		err := rows.Scan(&jobStatus.JobID, &jobStatus.Cutting, &jobStatus.Welding, &jobStatus.QualityCheck, &jobStatus.Packaging, &jobStatus.Dispatch)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jobStatuses = append(jobStatuses, jobStatus)
	}

	if len(jobStatuses) == 0 {
		jobStatuses = []jobsSchema.GetJobStatusResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobStatuses)
}
