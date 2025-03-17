package store

import (
	"database/sql"
	"errors"
	"github.com/AmoabaKelvin/temp-mail/api/internal/validator"
	"time"
)

type Message struct {
	ID          int64     `json:"id"`
	FromAddress string    `json:"from_address"`
	ToAddressID int64     `json:"to_address_id"`
	Subject     string    `json:"subject"`
	Body        string    `json:"body"`
	ReceivedAt  time.Time `json:"received_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageStore struct {
	db *sql.DB
}

func ValidateMessage(v *validator.Validator, message *Message) {
	v.Check(message.FromAddress != "", "from_address", "must be provided")
	v.Check(message.ToAddressID > 0, "to_address_id", "must be provided")
}

func (m MessageStore) Insert(message *Message) error {
	query := `
        INSERT INTO messages (from_address, to_address_id, subject, body, received_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	args := []any{message.FromAddress, message.ToAddressID, message.Subject, message.Body, message.ReceivedAt, time.Now(), time.Now()}
	return m.db.QueryRow(query, args...).Scan(&message.ID)
}

func (m MessageStore) Get(id int64) (*Message, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, from_address, to_address_id, subject, body, received_at, created_at, updated_at
        FROM messages
        WHERE id = $1 AND deleted_at IS NULL`

	var message Message
	err := m.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.FromAddress,
		&message.ToAddressID,
		&message.Subject,
		&message.Body,
		&message.ReceivedAt,
		&message.CreatedAt,
		&message.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &message, nil
}

func (m MessageStore) GetByFromAddress(fromAddress string) ([]*Message, error) {
	if fromAddress == "" {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, from_address, to_address_id, subject, body, received_at, created_at, updated_at
        FROM messages
        WHERE from_address = $1 AND deleted_at IS NULL`

	rows, err := m.db.Query(query, fromAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		var message Message

		err := rows.Scan(
			&message.ID,
			&message.FromAddress,
			&message.ToAddressID,
			&message.Subject,
			&message.Body,
			&message.ReceivedAt,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, ErrRecordNotFound
	}

	return messages, nil
}

func (m MessageStore) GetByToAddress(toAddressID int64) ([]*Message, error) {
	if toAddressID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, from_address, to_address_id, subject, body, received_at, created_at, updated_at
        FROM messages
        WHERE to_address_id = $1 AND deleted_at IS NULL`

	rows, err := m.db.Query(query, toAddressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		var message Message

		err := rows.Scan(
			&message.ID,
			&message.FromAddress,
			&message.ToAddressID,
			&message.Subject,
			&message.Body,
			&message.ReceivedAt,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, ErrRecordNotFound
	}

	return messages, nil
}

func (m MessageStore) Update(message *Message) error {
	query := `
        UPDATE messages
        SET from_address = $1, to_address_id = $2, subject = $3, body = $4, received_at = $5, updated_at = $6
        WHERE id = $7 AND deleted_at IS NULL`

	args := []any{
		message.FromAddress,
		message.ToAddressID,
		message.Subject,
		message.Body,
		message.ReceivedAt,
		time.Now(),
		message.ID,
	}

	result, err := m.db.Exec(query, args...)
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

func (m MessageStore) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM messages WHERE id = $1`

	result, err := m.db.Exec(query, id)
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

func (m MessageStore) GetAll(toAddressID int64) ([]*Message, error) {
	if toAddressID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, from_address, to_address_id, subject, body, received_at, created_at, updated_at
        FROM messages
        WHERE to_address_id = $1`

	rows, err := m.db.Query(query, toAddressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		var message Message

		err := rows.Scan(
			&message.ID,
			&message.FromAddress,
			&message.ToAddressID,
			&message.Subject,
			&message.Body,
			&message.ReceivedAt,
			&message.CreatedAt,
			&message.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, ErrRecordNotFound
	}

	return messages, nil
}
