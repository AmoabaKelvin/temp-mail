package config

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	TempmailDomains   []string
	ExpirationEnabled bool
	ExpireAfter       time.Duration
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	tempMailDomains := os.Getenv("TEMPMAIL_DOMAINS")
	if tempMailDomains == "" {
		return nil, errors.New("TEMPMAIL_DOMAINS is not set")
	}

	expirationEnabled := os.Getenv("EXPIRATION_ENABLED")
	if expirationEnabled == "" {
		expirationEnabled = "true"
	}

	expireAfter := os.Getenv("EXPIRE_AFTER")
	if expireAfter == "" {
		// set it to the default of 24 hours
		expireAfter = "24h"
	}

	expireAfterDuration, err := time.ParseDuration(expireAfter)
	if err != nil {
		return nil, err
	}

	return &Config{
		TempmailDomains:   strings.Split(tempMailDomains, ","),
		ExpirationEnabled: expirationEnabled == "true",
		ExpireAfter:       expireAfterDuration,
		DatabaseURL:       databaseURL,
	}, nil
}
