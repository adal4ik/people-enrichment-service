package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Router(handlers Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Post("/person", handlers.PersonHandler.CreatePerson)
	r.Get("/persons", handlers.PersonHandler.GetPersons)
	r.Get("/person/{id}", handlers.PersonHandler.GetPerson)
	r.Delete("/person/{id}", handlers.PersonHandler.DeletePerson)
	r.Put("/person/{id}", handlers.PersonHandler.UpdatePerson)

	return r
}
