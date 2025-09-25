package handlers

import (
	"database/sql"
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
	CreatedByID int    `json:"created_by_id"`
}

func getAllJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Getting Jobs-------")

	usernameCtx := r.Context().Value("username")
	roleCtx := r.Context().Value("role")

	if usernameCtx == nil || roleCtx == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := usernameCtx.(string)
	role := roleCtx.(string)

	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var rows *sql.Rows
	if role == "admin" {
		rows, err = db.Query("SELECT id, title, description, created_at, updated_at, employee_id, created_by_id FROM jobs")
	} else if role == "employee" {
		// get employee ID
		var employee_id int
		err = db.QueryRow("SELECT id FROM users WHERE username = ? LIMIT 1", username).Scan(&employee_id)
		if err != nil {
			http.Error(w, "invalid username", http.StatusUnauthorized)
			return
		}
		rows, err = db.Query("SELECT id, title, description, created_at, updated_at, employee_id, created_by_id FROM jobs WHERE employee_id = ?", employee_id)
	} else {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []JobResponse
	for rows.Next() {
		var job JobResponse
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.CreatedAt, &job.UpdatedAt, &job.EmployeeID, &job.CreatedByID); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		jobs = []JobResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
