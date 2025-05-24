package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
)

type Address struct {
	ID        int64      `json:"-"`
	Email     string     `json:"email"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	ReceivedMessages []Message `json:"-"`
}

type AddressStore struct {
	db *db.DB
}

func NewAddressStore(db *db.DB) *AddressStore {
	return &AddressStore{db: db}
}

func (s *AddressStore) Create(ctx context.Context, address *Address) error {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `INSERT INTO addresses (email, expires_at) VALUES ($1, $2) RETURNING id`
	return s.db.QueryRowContext(ctx, query, address.Email, address.ExpiresAt).Scan(&address.ID)
}

func (s *AddressStore) Get(ctx context.Context, email string) (*Address, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `SELECT id, email, expires_at FROM addresses WHERE email = $1`
	address := &Address{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(&address.ID, &address.Email, &address.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return address, err
}
