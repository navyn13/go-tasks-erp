package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type JobResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	EmployeeID  int    `json:"employee_id"`
}

func getAllJobsForEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Getting Jobs for Employee-------")
	EmployeeUsernameCtx := r.Context().Value("username")
	if EmployeeUsernameCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	employee_username := EmployeeUsernameCtx.(string)
	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var employee_id int
	err = db.QueryRow("SELECT id FROM users WHERE username = ? LIMIT 1", employee_username).Scan(&employee_id)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	rows, err := db.Query("SELECT id, title, description , created_at, updated_at, employee_id FROM jobs WHERE employee_id = ?", employee_id)
	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []JobResponse
	for rows.Next() {
		var job JobResponse
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.CreatedAt, &job.UpdatedAt, &job.EmployeeID)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, job)
	}

	// if no jobs found, return empty array
	if len(jobs) == 0 {
		jobs = []JobResponse{}
	}

	response := jobs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
