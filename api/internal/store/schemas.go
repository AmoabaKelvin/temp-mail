package store

import (
	"time"
)

type AddressSchema struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}

type MessageSchema struct {
	ID          int64     `json:"id"`
	FromAddress string    `json:"from_address"`
	ToAddressID int64     `json:"to_address_id"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	ReceivedAt  time.Time `json:"received_at"`
}
