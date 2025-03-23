package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/repository"
	"github.com/AmoabaKelvin/temp-mail/pkg/config"
	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Handler struct {
	repo   *repository.Repository
	config *config.Config
}

func New(repo *repository.Repository, config *config.Config) *Handler {
	return &Handler{repo: repo, config: config}
}

func (h *Handler) generateAddress() models.Address {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}

	domain := h.config.TempmailDomains[rand.Intn(len(h.config.TempmailDomains))]
	return models.Address{
		Email:     fmt.Sprintf("%s@%s", id, domain),
		ExpiresAt: time.Now().Add(h.config.ExpireAfter),
	}
}

func (h *Handler) GenerateAddress(w http.ResponseWriter, r *http.Request) {
	address := h.generateAddress()

	if err := h.repo.InsertAddress(&address); err != nil {
		http.Error(w, "Failed to insert address", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(address)
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if message.FromAddress == "" || message.ToAddressID == 0 || message.ReceivedAt.IsZero() {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	isExpired, err := h.repo.IsAddressExpired(message.ToAddressID)
	if err != nil {
		http.Error(w, "Failed to check address expiry", http.StatusInternalServerError)
		return
	}

	if isExpired {
		http.Error(w, "Address expired", http.StatusBadRequest)
		return
	}

	if err := h.repo.InsertMessage(&message); err != nil {
		http.Error(w, "Failed to insert message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	address, err := h.repo.GetAddressByEmail(email)
	if err == repository.ErrRecordNotFound {
		http.Error(w, "Recipient not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	isExpired, err := h.repo.IsAddressExpired(address.ID)
	if err != nil {
		http.Error(w, "Failed to check address expiry", http.StatusInternalServerError)
		return
	}

	if isExpired {
		http.Error(w, "Address expired", http.StatusBadRequest)
		return
	}

	messages, err := h.repo.GetMessagesByRecipient(address.ID)
	if err == repository.ErrRecordNotFound {
		json.NewEncoder(w).Encode([]models.Message{})
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	err = h.repo.DeleteMessage(uint(id))
	if err == repository.ErrRecordNotFound {
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
