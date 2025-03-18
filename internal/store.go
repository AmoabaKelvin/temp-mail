package internal

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Address struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	ExpiresAt string `json:"expires_at"`
}

type Message struct {
	ID          int64  `json:"id"`
	FromAddress string `json:"from_address"`
	ToAddressID int64  `json:"to_address_id"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ReceivedAt  string `json:"received_at"`
}

func InsertAddress(db *sql.DB, address *Address) error {
	query := `INSERT INTO addresses (email, expires_at) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, address.Email, address.ExpiresAt).Scan(&address.ID)
}

func GetAddressByEmail(db *sql.DB, email string) (*Address, error) {
	address := &Address{}
	query := `SELECT id, email, expires_at FROM addresses WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&address.ID, &address.Email, &address.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	return address, err
}

func InsertMessage(db *sql.DB, message *Message) error {
	query := `INSERT INTO messages (from_address, to_address_id, subject, body, received_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.QueryRow(query, message.FromAddress, message.ToAddressID, message.Subject, message.Body, message.ReceivedAt).Scan(&message.ID)
}

func GetMessagesByRecipient(db *sql.DB, toAddressID int64) ([]Message, error) {
	query := `SELECT id, from_address, to_address_id, subject, body, received_at FROM messages WHERE to_address_id = $1`
	rows, err := db.Query(query, toAddressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.FromAddress, &msg.ToAddressID, &msg.Subject, &msg.Body, &msg.ReceivedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return nil, ErrRecordNotFound
	}
	return messages, nil
}

func DeleteMessage(db *sql.DB, id int64) error {
	query := `DELETE FROM messages WHERE id = $1`
	res, err := db.Exec(query, id)
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
