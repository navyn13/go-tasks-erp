package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/navyn13/go-tasks-erp/internal/db"
	"github.com/navyn13/go-tasks-erp/internal/models/jobsSchema"
)

func UpdateJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Updating Job-------")
	userIdCtx := r.Context().Value("id")
	userID := userIdCtx.(int)

	var req jobsSchema.UpdateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.JobID == 0 {
		http.Error(w, "Title, Description and JobID are required", http.StatusBadRequest)
		return
	}

	dbConn, err := db.GetDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	result, err := dbConn.Exec("UPDATE jobs SET title = ?, description = ?, employee_id = ? WHERE id = ?",
		req.Title, req.Description, userID, req.JobID)
	if err != nil {
		http.Error(w, "Unable to update job: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update result: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "No job found with given ID", http.StatusNotFound)
		return
	}

	response := jobsSchema.UpdateJobResponse{Message: "Job Updated Successfully"}
	fmt.Println(response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
