package store

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	Message interface {
		Insert(message *Message) error
		Get(id int64) (*Message, error)
		GetByFromAddress(fromAddress string) ([]*Message, error)
		GetByToAddress(toAddressID int64) ([]*Message, error)
		Delete(id int64) error
		GetAll(toAddressID int64) ([]*Message, error)
		Update(message *Message) error
	}

	Address interface {
		Insert(address *Address) error
		Get(id int64) (*Address, error)
		GetByEmail(email string) (*Address, error)
		Delete(id int64) error
		GetAll() ([]*Address, error)
		Update(address *Address) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Message: &MessageStore{db},
		Address: &AddressStore{db},
	}
}
