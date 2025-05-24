package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/store"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *application) newRandomAddress() store.Address {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}

	domain := app.config.tempMail.domains[rand.Intn(len(app.config.tempMail.domains))]
	duration, err := time.ParseDuration(app.config.tempMail.expireAfter)
	if err != nil {
		panic(err)
	}

	return store.Address{
		Email:     fmt.Sprintf("%s@%s", id, domain),
		ExpiresAt: time.Now().Add(duration),
	}
}

func (app *application) generateAddress(w http.ResponseWriter, r *http.Request) {
	address := app.newRandomAddress()

	if err := app.store.Addresses.Create(r.Context(), &address); err != nil {
		http.Error(w, "Failed to insert address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}
