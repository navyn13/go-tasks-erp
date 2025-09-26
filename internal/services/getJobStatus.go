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

	var req jobsSchema.GetJobProcessStatusRequest
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

	var jobStatus jobsSchema.GetJobProcessStatusResponse
	var row *sql.Row
	if role == "admin" {
		row = dbConn.QueryRow(
			"SELECT jobid, process_name, status, started_at, completed_at FROM jobStatus WHERE jobid = ? AND process_name = ?",
			req.JobID, req.Process)
	} else if role == "employee" {
		updateProcessQuery := `
		SELECT js.jobid, js.process_name, js.status, js.started_at, js.completed_at
		FROM jobStatus js
		JOIN jobs j ON j.id = js.jobid
		WHERE js.jobid = ? AND j.employee_id = ?`

		row = dbConn.QueryRow(updateProcessQuery, req.JobID, userId)
	} else {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}
	err = row.Scan(&jobStatus.JobID, &jobStatus.Process, &jobStatus.Status, &jobStatus.StartedAt, &jobStatus.CompletedAt)
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
