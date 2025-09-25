package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/navyn13/go-tasks-erp/internal/middleware"
)

func Handlers(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello this is handlers- API is working :)"))
	})
	r.Post("/signup", signup)
	r.Post("/login", login)

	r.Route("/admin", func(router chi.Router) {
		router.Use(middleware.AdminOnly)       // admin ke liye routes
		router.Get("/jobs", getAllJobs)        // get all jobs jitni bhi hae
		router.Get("/jobStatus", getJobStatus) // get job status for admin
		router.Post("/jobs", createjob)        // create a job
		router.Put("/jobs", updateJob)         // update a job
		router.Delete("/jobs", deletejob)      // delete a job

	})
	r.Route("/employee", func(router chi.Router) {
		router.Use(middleware.EmployeeOnly)             // employee routes
		router.Get("/jobs", getAllJobs)                 // get all jobs for employee
		router.Get("/jobstatus", getJobStatus)          // get job status for employee
		router.Put("/updatejobstatus", updateJobStatus) // update job status for employee
	})

}
