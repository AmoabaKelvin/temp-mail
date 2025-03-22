package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AmoabaKelvin/temp-mail/pkg/config"
	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func generateAddress(config *config.Config) models.Address {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}

	// select a random domain from the config
	domain := config.TempmailDomains[rand.Intn(len(config.TempmailDomains))]

	return models.Address{
		Email:     fmt.Sprintf("%s@%s", id, domain),
		ExpiresAt: time.Now().Add(config.ExpireAfter),
	}
}

func GenerateAddressHandler(db *sql.DB, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := generateAddress(config)

		if err := InsertAddress(db, &address); err != nil {
			http.Error(w, "Failed to insert address", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(address)
	}
}

func CreateMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message models.Message

		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if message.FromAddress == "" || message.ToAddressID == 0 || message.ReceivedAt.IsZero() {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if err := InsertMessage(db, &message); err != nil {
			http.Error(w, "Failed to insert message", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(message)
	}
}

func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email parameter is required", http.StatusBadRequest)
			return
		}

		address, err := GetAddressByEmail(db, email)
		if err == ErrRecordNotFound {
			http.Error(w, "Recipient not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		messages, err := GetMessagesByRecipient(db, address.ID)
		if err == ErrRecordNotFound {
			json.NewEncoder(w).Encode([]models.Message{})
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(messages)
	}
}

func DeleteMessageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = DeleteMessage(db, uint(id))
		if err == ErrRecordNotFound {
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
}
