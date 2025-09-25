package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type JobStatusResponse struct {
	Job_id        int    `json:"job_id"`
	Cutting       string `json:"cutting"`
	Welding       string `json:"welding"`
	Quality_check string `json:"quality_check"`
	Packaging     string `json:"packaging"`
	Dispatch      string `json:"dispatch"`
}
type JobStatusUpdateRequest struct {
	JobID   int    `json:"job_id"`
	Process string `json:"process"`
	Status  string `json:"status"`
}

func updateJobStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------Updating Job Status--------")
	db, err := getDB()
	if err != nil {
		http.Error(w, "DB connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req JobStatusUpdateRequest
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
	//query building
	updateProcessQuery := fmt.Sprintf("UPDATE jobStatus SET %s = ? WHERE job_id = ?", req.Process)

	_, err = db.Exec(updateProcessQuery, req.Status, req.JobID)
	if err != nil {
		http.Error(w, "DB update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job status updated successfully"))
}
