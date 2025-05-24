package main

import (
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	config *config
	store  *store.Storage
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
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3000/*", "https://www.is-temp.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/addresses", app.generateAddress)
		r.Get("/messages", app.getMessages)
		r.Delete("/messages/{id}", app.deleteMessage)
		r.Put("/messages/{id}/read", app.updateMessageReadAt)
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
