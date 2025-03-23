package models

import (
	"fmt"
	"log"
	"time"

	"github.com/AmoabaKelvin/temp-mail/pkg/config"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Address represents an email address with expiration
type Address struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	ExpiresAt time.Time      `gorm:"index" json:"expires_at"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	ReceivedMessages []Message `gorm:"foreignKey:ToAddressID" json:"-"`
}

// Message represents a simple email message
type Message struct {
	ID          uint           `gorm:"primaryKey"`
	FromAddress string         `gorm:"type:varchar(255);not null"`
	ToAddressID uint           `gorm:"index" json:"-"`
	ToAddress   Address        `gorm:"foreignKey:ToAddressID" json:"-"`
	Headers     datatypes.JSON `gorm:"type:jsonb"`
	Subject     string         `gorm:"type:varchar(255)"`
	Body        string         `gorm:"type:text"`
	ReceivedAt  time.Time      `gorm:"index"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func SetupDatabase(config *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Address{}, &Message{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed")
	return db
}
