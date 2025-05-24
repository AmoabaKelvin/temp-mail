package store

import (
	"context"
	"errors"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
)

var (
	ErrNotFound          = errors.New("record not found")
	QueryDurationTimeout = 5 * time.Second // default timeout for queries
)

type Storage struct {
	Messages interface {
		Get(context.Context, int64) ([]Message, error)
		Delete(context.Context, int64) error
		SetReadAt(context.Context, int64, *time.Time) error
		Create(context.Context, *Message) error
	}
	Addresses interface {
		Create(context.Context, *Address) error
		Get(context.Context, string) (*Address, error)
	}
}

func NewStorage(db *db.DB) *Storage {
	return &Storage{
		Messages:  NewMessageStore(db),
		Addresses: NewAddressStore(db),
	}
}
