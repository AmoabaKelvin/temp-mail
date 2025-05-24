package main

import (
	"log"
	"os"
	"strings"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
	"github.com/AmoabaKelvin/temp-mail/internal/store"
)

func main() {
	config := &config{
		addr: os.Getenv("ADDR"),
		db: &dbConfig{
			addr: os.Getenv("DATABASE_URL"),
		},
		tempMail: &tempMailConfig{
			domains:           strings.Split(os.Getenv("TEMPMAIL_DOMAINS"), ","),
			expireAfter:       os.Getenv("EXPIRE_AFTER"),
			expirationEnabled: os.Getenv("EXPIRATION_ENABLED"),
		},
	}

	db, err := db.New(config.db.addr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	store := store.NewStorage(db)
	app := &application{
		config: config,
		store:  store,
	}

	routes := app.mount()

	if err := app.run(routes); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
