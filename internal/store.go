package internal

import (
	"database/sql"
	"errors"

	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
)

var ErrRecordNotFound = errors.New("record not found")

func InsertAddress(db *sql.DB, address *models.Address) error {
	query := `INSERT INTO addresses (email, expires_at) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, address.Email, address.ExpiresAt).Scan(&address.ID)
}

func GetAddressByEmail(db *sql.DB, email string) (*models.Address, error) {
	address := &models.Address{}
	query := `SELECT id, email, expires_at FROM addresses WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&address.ID, &address.Email, &address.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	return address, err
}

func InsertMessage(db *sql.DB, message *models.Message) error {
	query := `INSERT INTO messages (from_address, to_address_id, subject, body, received_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.QueryRow(query, message.FromAddress, message.ToAddressID, message.Subject, message.Body, message.ReceivedAt).Scan(&message.ID)
}

func GetMessagesByRecipient(db *sql.DB, toAddressID uint) ([]models.Message, error) {
	query := `SELECT id, from_address, to_address_id, subject, body, received_at FROM messages WHERE to_address_id = $1`
	rows, err := db.Query(query, toAddressID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
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

func DeleteMessage(db *sql.DB, id uint) error {
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
