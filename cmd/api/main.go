package main

import (
	"log"
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/internal"
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

	db := internal.ConnectDB(config.DatabaseURL)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/v1/generate_address", internal.GenerateAddressHandler(db))
	r.Post("/v1/create_message", internal.CreateMessageHandler(db))
	r.Get("/v1/get_messages", internal.GetMessagesHandler(db))
	r.Delete("/v1/delete_message", internal.DeleteMessageHandler(db))

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
