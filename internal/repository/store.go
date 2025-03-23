package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/database"
	"github.com/AmoabaKelvin/temp-mail/pkg/config"
	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
)

var ErrRecordNotFound = errors.New("record not found")

type Repository struct {
	db     *database.DB
	config *config.Config
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertAddress(address *models.Address) error {
	query := `INSERT INTO addresses (email, expires_at) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, address.Email, address.ExpiresAt).Scan(&address.ID)
}

func (r *Repository) GetAddressByEmail(email string) (*models.Address, error) {
	address := &models.Address{}
	query := `SELECT id, email, expires_at FROM addresses WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&address.ID, &address.Email, &address.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	return address, err
}

func (r *Repository) InsertMessage(message *models.Message) error {
	query := `INSERT INTO messages (from_address, to_address_id, subject, body, received_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(
		query,
		message.FromAddress,
		message.ToAddressID,
		message.Subject,
		message.Body,
		message.ReceivedAt,
	).Scan(&message.ID)
}

func (r *Repository) GetMessagesByRecipient(toAddressID uint) ([]models.Message, error) {
	query := `SELECT id, from_address, to_address_id, subject, body, received_at
              FROM messages WHERE to_address_id = $1`
	rows, err := r.db.Query(query, toAddressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.FromAddress,
			&msg.ToAddressID,
			&msg.Subject,
			&msg.Body,
			&msg.ReceivedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return nil, ErrRecordNotFound
	}
	return messages, nil
}

func (r *Repository) DeleteMessage(id uint) error {
	query := `DELETE FROM messages WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *Repository) IsAddressExpired(id uint) (bool, error) {
	expirationEnabled := r.config.ExpirationEnabled
	if !expirationEnabled {
		return false, nil
	}

	query := `SELECT expires_at FROM addresses WHERE id = $1`
	var expiresAt int64
	err := r.db.QueryRow(query, id).Scan(&expiresAt)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("database error checking expiration: %v", err)
	}

	return expiresAt < time.Now().Unix(), nil
}
