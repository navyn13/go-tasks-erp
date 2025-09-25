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

	// check username in context
	employeeUsernameCtx := r.Context().Value("username")
	if employeeUsernameCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	employee_username := employeeUsernameCtx.(string)

	// get db
	db, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// get employee id
	var employee_id int
	err = db.QueryRow("SELECT id FROM users WHERE username = ? LIMIT 1", employee_username).Scan(&employee_id)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusUnauthorized)
		return
	}

	// decode request body
	var req jobsSchema.UpdateJobStatusRequest
	allowedProcesses := map[string]bool{
		"cutting":       true,
		"welding":       true,
		"quality_check": true,
		"packaging":     true,
		"dispatch":      true,
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Process == "" || req.Status == "" || req.JobID == 0 || !allowedProcesses[req.Process] {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// query with ownership check (employee must own the job)
	updateProcessQuery := fmt.Sprintf(`
		UPDATE jobStatus js
		JOIN jobs j ON j.id = js.job_id
		SET js.%s = ?
		WHERE js.job_id = ? AND j.employee_id = ?`, req.Process)

	res, err := db.Exec(updateProcessQuery, req.Status, req.JobID, employee_id)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// check if any row was updated
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Not allowed: either job not found or not assigned to you", http.StatusForbidden)
		return
	}

	// response
	response := jobsSchema.UpdateJobStatusResponse{Message: "Job status updated successfully"}
	fmt.Println(response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
