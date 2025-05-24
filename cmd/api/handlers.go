package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/store"
	"github.com/go-chi/chi/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *application) generateAddress() store.Address {
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

func (app *application) GenerateAddress(w http.ResponseWriter, r *http.Request) {
	address := app.generateAddress()

	if err := app.store.Addresses.Create(r.Context(), &address); err != nil {
		http.Error(w, "Failed to insert address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

func (app *application) GetMessages(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	address, err := app.store.Addresses.Get(r.Context(), email)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	messages, err := app.store.Messages.Get(r.Context(), address.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (app *application) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	if err := app.store.Messages.Delete(r.Context(), id); err != nil {
		switch err {
		case store.ErrNotFound:
			http.Error(w, "Message not found", http.StatusNotFound)
		default:
			http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "Message deleted successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (app *application) UpdateMessageReadAt(w http.ResponseWriter, r *http.Request) {
	readAt := time.Now()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = app.store.Messages.SetReadAt(r.Context(), id, &readAt)
	if err != nil {
		http.Error(w, "Failed to update message read status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
