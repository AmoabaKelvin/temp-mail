package main

import (
	"log"
	"os"

	"github.com/AmoabaKelvin/temp-mail/internal/mailserver"
)

func main() {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalf("DATABASE_URL is not set")
	}
	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		log.Fatalf("SMTP_PORT is not set")
	}

	if err := mailserver.Start(databaseUrl, smtpPort); err != nil {
		log.Fatalf("Failed to start mail server: %v", err)
	}
}
