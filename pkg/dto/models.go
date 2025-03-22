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
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	ReceivedMessages []Message `gorm:"foreignKey:ToAddressID"`
}

// Message represents a simple email message
type Message struct {
	ID          uint           `gorm:"primaryKey"`
	FromAddress string         `gorm:"type:varchar(255);not null"`
	ToAddressID uint           `gorm:"index"`
	ToAddress   Address        `gorm:"foreignKey:ToAddressID"`
	Headers     datatypes.JSON `gorm:"type:jsonb"`
	Subject     string         `gorm:"type:varchar(255)"`
	Body        string         `gorm:"type:text"`
	ReceivedAt  time.Time      `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
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
