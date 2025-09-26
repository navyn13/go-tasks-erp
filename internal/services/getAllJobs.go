package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/navyn13/go-tasks-erp/internal/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Getting Jobs-------")

	roleCtx := r.Context().Value("role")
	userIdCtx := r.Context().Value("id")

	userId := userIdCtx.(int)
	role := roleCtx.(string)

	db, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var rows *sql.Rows
	if role == "admin" {
		rows, err = db.Query("SELECT id, title, description, created_at, updated_at, employee_id, created_by_id FROM jobs")
	} else if role == "employee" {
		rows, err = db.Query("SELECT id, title, description, created_at, updated_at, employee_id, created_by_id FROM jobs WHERE employee_id = ?", userId)
	} else {
		http.Error(w, "Invalid role", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, "DB query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	response := jobsSchema.GetAllJobsResponse{
		Jobs: []jobsSchema.Job{},
	}

	for rows.Next() {
		var job jobsSchema.Job
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.CreatedAt, &job.UpdatedAt, &job.EmployeeID, &job.CreatedByID); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response.Jobs = append(response.Jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
