package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
)

func Printer(s string) {
	fmt.Println(s)
}

func Handlers(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello this is handlers- API is working :)"))
	})

	r.Post("/signup", signup)
	r.Post("/login", login)

}
