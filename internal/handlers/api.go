package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/navyn13/go-tasks-erp/internal/middleware"
	"github.com/navyn13/go-tasks-erp/internal/services"
	"github.com/navyn13/go-tasks-erp/internal/tools"
)

func Handlers(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello this is handlers- API is working :)"))
	})
	r.Post("/signup", tools.Signup)
	r.Post("/login", tools.Login)

	r.Route("/admin", func(router chi.Router) {
		router.Use(middleware.AdminOnly)                // admin ke liye routes
		router.Get("/jobs", services.GetAllJobs)        // get all jobs jitni bhi hae
		router.Post("/jobs", services.CreateJob)        // create a job
		router.Put("/jobs", services.UpdateJob)         // update a job
		router.Delete("/jobs", services.DeleteJob)      // delete a job
		router.Get("/jobStatus", services.GetJobStatus) // get job status for admin

	})
	r.Route("/employee", func(router chi.Router) {
		router.Use(middleware.EmployeeOnly)                      // employee routes
		router.Get("/jobs", services.GetAllJobs)                 // get all jobs for employee
		router.Get("/jobstatus", services.GetJobStatus)          // get job status for employee (no employee should get the data of other employees jobs)
		router.Put("/updatejobstatus", services.UpdateJobStatus) // update job status for employee (no employee should update the data of other employees jobs)
	})

}
