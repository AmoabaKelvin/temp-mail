package models

import (
	"time"
)

type Address struct {
	ID        uint       `json:"-"`
	Email     string     `json:"email"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	ReceivedMessages []Message `json:"-"`
}

// rather than having a boolean flag for the read status, we will rather use a timestamp. with this,
// we can not only know if the message has been read, but also the exact time it was read which can help
// implement time-based email filters
type Message struct {
	ID          uint       `json:"id"`
	FromAddress string     `json:"from_address"`
	ToAddressID uint       `json:"-"`
	ToAddress   Address    `json:"-"`
	Headers     []byte     `json:"headers"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	ReceivedAt  time.Time  `json:"received_at"`
	ReadAt      time.Time  `json:"read_at"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}
