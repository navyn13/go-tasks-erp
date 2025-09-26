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

	userIdCtx, ok := r.Context().Value("id").(int)
	if !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusUnauthorized)
		return
	}
	userId := userIdCtx

	// Safely get the role
	roleCtx, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "Invalid or missing role in context", http.StatusUnauthorized)
		return
	}
	role := roleCtx
	fmt.Println("userid: ", userId)
	fmt.Println("role: ", role)

	var req jobsSchema.GetJobStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var jobStatus jobsSchema.GetJobStatusResponse
	var row *sql.Row
	if role == "admin" {
		row = dbConn.QueryRow(
			"SELECT job_id, cutting, welding, quality_check, packaging, dispatch FROM jobStatus WHERE job_id = ?",
			req.JobID)
	} else if role == "employee" {
		updateProcessQuery := `
		SELECT js.job_id, js.cutting, js.welding, js.quality_check, js.packaging, js.dispatch
		FROM jobStatus js
		JOIN jobs j ON j.id = js.job_id
		WHERE js.job_id = ? AND j.employee_id = ?`

		row = dbConn.QueryRow(updateProcessQuery, req.JobID, userId)
	} else {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}
	err = row.Scan(&jobStatus.JobID, &jobStatus.Cutting, &jobStatus.Welding, &jobStatus.QualityCheck, &jobStatus.Packaging, &jobStatus.Dispatch)
	if err != nil {
		if err == sql.ErrNoRows {

			http.Error(w, "Job status not found", http.StatusNotFound)
		} else {
			http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobStatus)
}
