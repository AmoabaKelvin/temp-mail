package main

import (
	"log"
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/internal/database"
	"github.com/AmoabaKelvin/temp-mail/internal/handlers"
	"github.com/AmoabaKelvin/temp-mail/internal/repository"
	"github.com/AmoabaKelvin/temp-mail/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config")
	}

	db, err := database.New(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repository := repository.New(db)
	handler := handlers.New(repository, config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/v1/addresses", handler.GenerateAddress)
	r.Get("/v1/messages", handler.GetMessages)
	r.Delete("/v1/messages/{id}", handler.DeleteMessage)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
