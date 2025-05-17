package main

import (
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	config *config
	store  *repository.Repository
}

type config struct {
	addr     string
	db       *dbConfig
	tempMail *tempMailConfig
}

type dbConfig struct {
	addr string
}

type tempMailConfig struct {
	domains           []string
	expireAfter       string
	expirationEnabled string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1/", func(r chi.Router) {
		r.Get("/addresses", app.GenerateAddress)
		r.Get("/messages", app.GenerateAddress)
		r.Get("/messages/{id}", app.GenerateAddress)
		r.Get("/messages/{id}/read", app.GenerateAddress)
	})

	return r
}

func (app *application) run(routes http.Handler) error {
	server := &http.Server{
		Addr:    app.config.addr,
		Handler: routes,
	}

	return server.ListenAndServe()
}
