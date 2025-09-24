package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/navyn13/go-tasks-erp/internal/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {

	var r *chi.Mux = chi.NewRouter()
	handlers.Handlers(r)

	fmt.Println("STARTING GO-TASKS-ERP API SERVICES")
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Error(err)
	}
}
