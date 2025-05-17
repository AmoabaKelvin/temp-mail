package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/repository"
	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
	"github.com/go-chi/chi/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *application) generateAddress() models.Address {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}

	domain := app.config.tempMail.domains[rand.Intn(len(app.config.tempMail.domains))]
	duration, err := time.ParseDuration(app.config.tempMail.expireAfter)
	if err != nil {
		panic(err)
	}

	return models.Address{
		Email:     fmt.Sprintf("%s@%s", id, domain),
		ExpiresAt: time.Now().Add(duration),
	}
}

func (app *application) GenerateAddress(w http.ResponseWriter, r *http.Request) {
	address := app.generateAddress()

	if err := app.store.InsertAddress(&address); err != nil {
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

	address, err := app.store.GetAddressByEmail(email)
	if errors.Is(err, repository.ErrRecordNotFound) {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	messages, err := app.store.GetMessagesByRecipient(address.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (app *application) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = app.store.DeleteMessage(uint(id))
	if errors.Is(err, repository.ErrRecordNotFound) {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Message deleted successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (app *application) UpdateMessageReadAt(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = app.store.UpdateMessageReadAt(uint(id), time.Now())
	if err != nil {
		http.Error(w, "Failed to update message read status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
