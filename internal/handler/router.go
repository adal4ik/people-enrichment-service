package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router(handlers Handler) http.Handler {
	r := chi.NewRouter()

	r.Post("/person", handlers.PersonHandler.CreatePerson)

	return r
}
