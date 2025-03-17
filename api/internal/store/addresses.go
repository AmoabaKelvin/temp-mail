package store

import (
	"database/sql"
	"errors"
	"github.com/AmoabaKelvin/temp-mail/api/internal/validator"
	"time"
)

type Address struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddressStore struct {
	db *sql.DB
}

func ValidateAddress(v *validator.Validator, address *Address) {
	v.Check(address.Email != "", "email", "must be provided")
	v.Check(address.ExpiresAt.After(time.Now()), "expires_at", "must be a date in the future")
}

func (a AddressStore) Insert(address *Address) error {
	query := `
        INSERT INTO addresses (email, expires_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	args := []any{address.Email, address.ExpiresAt, time.Now(), time.Now()}
	return a.db.QueryRow(query, args...).Scan(&address.ID)
}

func (a AddressStore) Get(id int64) (*Address, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, email, expires_at, created_at, updated_at
        FROM addresses
        WHERE id = $1 AND deleted_at IS NULL`

	var address Address
	err := a.db.QueryRow(query, id).Scan(
		&address.ID,
		&address.Email,
		&address.ExpiresAt,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &address, nil
}

func (a AddressStore) GetByEmail(email string) (*Address, error) {
	if email == "" {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, email, expires_at, created_at, updated_at
        FROM addresses
        WHERE email = $1 AND deleted_at IS NULL`

	var address Address
	err := a.db.QueryRow(query, email).Scan(
		&address.ID,
		&address.Email,
		&address.ExpiresAt,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &address, nil
}

func (a AddressStore) Update(address *Address) error {
	query := `
        UPDATE addresses
        SET email = $1, expires_at = $2, updated_at = $3
        WHERE id = $4 AND deleted_at IS NULL`

	args := []any{
		address.Email,
		address.ExpiresAt,
		time.Now(),
		address.ID,
	}

	result, err := a.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (a AddressStore) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// Soft delete - update deleted_at field
	query := `
        UPDATE addresses
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL`

	result, err := a.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (a AddressStore) GetAll() ([]*Address, error) {
	query := `
        SELECT id, email, expires_at, created_at, updated_at
        FROM addresses
        WHERE deleted_at IS NULL`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := []*Address{}

	for rows.Next() {
		var address Address

		err := rows.Scan(
			&address.ID,
			&address.Email,
			&address.ExpiresAt,
			&address.CreatedAt,
			&address.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		addresses = append(addresses, &address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
