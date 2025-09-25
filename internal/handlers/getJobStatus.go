package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type JobStatus struct {
	JobID        int    `json:"job_id"`
	Cutting      string `json:"cutting"`
	Welding      string `json:"welding"`
	QualityCheck string `json:"quality_check"`
	Packaging    string `json:"packaging"`
	Dispatch     string `json:"dispatch"`
}
type JobStatusRequest struct {
	JobID int `json:"job_id"`
}

func getJobStatus(w http.ResponseWriter, r *http.Request) {
	// Placeholder implementation
	fmt.Println("--------Getting Job Status-------")

	var req JobStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	db, err := getDB()
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

	var jobStatuses []JobStatus
	for rows.Next() {
		var jobStatus JobStatus
		err := rows.Scan(&jobStatus.JobID, &jobStatus.Cutting, &jobStatus.Welding, &jobStatus.QualityCheck, &jobStatus.Packaging, &jobStatus.Dispatch)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jobStatuses = append(jobStatuses, jobStatus)
	}

	if len(jobStatuses) == 0 {
		jobStatuses = []JobStatus{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobStatuses)
}
