package jobsSchema

type CreateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EmployeeID  int    `json:"employee_id"`
}
type CreateJobResponse struct {
	Message string `json:"message"`
}

type DeleteJobRequest struct {
	JobID int `json:"job_id"`
}
type DeleteJobResponse struct {
	Message string `json:"message"`
}

type GetAllJobsResponse struct {
	Jobs []Job `json:"jobs"`
}
type Job struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	EmployeeID  int    `json:"employee_id"`
	CreatedByID int    `json:"created_by_id"`
}

type GetJobStatusResponse struct {
	JobID        int    `json:"job_id"`
	Cutting      string `json:"cutting"`
	Welding      string `json:"welding"`
	QualityCheck string `json:"quality_check"`
	Packaging    string `json:"packaging"`
	Dispatch     string `json:"dispatch"`
}
type GetJobStatusRequest struct {
	JobID int `json:"job_id"`
}

type UpdateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	JobID       int    `json:"job_id"`
}
type UpdateJobResponse struct {
	Message string `json:"message"`
}

type UpdateJobStatusResponse struct {
	Message string `json:"message"`
}
type UpdateJobStatusRequest struct {
	JobID   int    `json:"job_id"`
	Process string `json:"process"`
	Status  string `json:"status"`
}
