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

	var req jobsSchema.DeleteJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.JobID == 0 {
		http.Error(w, "JobID is required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()
	// transaction- if fails revert all
	tx, err := dbConn.Begin()
	if err != nil {
		http.Error(w, "DB transaction error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM jobStatus WHERE job_id = ?", req.JobID)
	if err != nil {
		http.Error(w, "JobStatus delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the job
	result, err := tx.Exec("DELETE FROM jobs WHERE id = ?", req.JobID)
	if err != nil {
		http.Error(w, "Job delete error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No job found with given ID", http.StatusNotFound)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "DB commit error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := jobsSchema.DeleteJobResponse{Message: "Job deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
