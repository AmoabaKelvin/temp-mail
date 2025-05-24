package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
)

type Message struct {
	ID          uint       `json:"id"`
	FromAddress string     `json:"from_address"`
	ToAddressID uint       `json:"-"`
	ToAddress   Address    `json:"-"`
	Headers     []byte     `json:"headers"`
	Subject     string     `json:"subject"`
	BodyHTML    *string    `json:"body_html"`
	BodyPlain   *string    `json:"body_plain"`
	ContentType string     `json:"content_type"`
	ReceivedAt  time.Time  `json:"received_at"`
	ReadAt      time.Time  `json:"read_at"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}

type MessageStore struct {
	db *db.DB
}

func NewMessageStore(db *db.DB) *MessageStore {
	return &MessageStore{db: db}
}

// Get all messages for a given address ID
func (s *MessageStore) Get(ctx context.Context, addressID int64) ([]Message, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `SELECT id, from_address, to_address_id, headers, subject, body_html, body_plain, content_type, received_at, read_at 
			FROM messages 
			WHERE to_address_id = $1 
			ORDER BY received_at DESC`
	messages := []Message{}
	rows, err := s.db.QueryContext(ctx, query, addressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.FromAddress, &message.ToAddressID, &message.Headers, &message.Subject, &message.BodyHTML, &message.BodyPlain, &message.ContentType, &message.ReceivedAt, &message.ReadAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (s *MessageStore) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `DELETE FROM messages WHERE id = $1`
	executionResult, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := executionResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *MessageStore) SetReadAt(ctx context.Context, id int64, readAt *time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `UPDATE messages SET read_at = $1 WHERE id = $2`
	_, err := s.db.ExecContext(ctx, query, readAt, id)
	return err
}

func (s *MessageStore) Create(ctx context.Context, message *Message) error {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `INSERT INTO messages (from_address, to_address_id, subject, body_html, body_plain, content_type, headers, received_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := s.db.QueryRowContext(ctx, query,
		message.FromAddress,
		message.ToAddressID,
		message.Subject,
		message.BodyHTML,
		message.BodyPlain,
		message.ContentType,
		message.Headers,
		message.ReceivedAt,
	).Scan(&message.ID)

	return err
}

// GetByID gets a single message by its ID
func (s *MessageStore) GetByID(ctx context.Context, messageID int64) (*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryDurationTimeout)
	defer cancel()

	query := `SELECT id, from_address, to_address_id, headers, subject, body_html, body_plain, content_type, received_at, read_at 
		FROM messages 
		WHERE id = $1`

	var message Message
	err := s.db.QueryRowContext(
		ctx,
		query,
		messageID,
	).Scan(
		&message.ID,
		&message.FromAddress,
		&message.ToAddressID,
		&message.Headers,
		&message.Subject,
		&message.BodyHTML,
		&message.BodyPlain,
		&message.ContentType,
		&message.ReceivedAt,
		&message.ReadAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &message, nil
}
