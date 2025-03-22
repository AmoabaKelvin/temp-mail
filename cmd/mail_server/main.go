package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	models "github.com/AmoabaKelvin/temp-mail/pkg/dto"
)

// todo
// check if the incoming email is valid and is in our database
// we can either check that in the rcpt method or in the mail method
// if we check in the rcpt method, that will be one round trip to the database
// but if we check in the mail method, we don't bother saving the email if the address is not valid
// parse email and store it in the database
//

// todo
//
// 1. read in the environment variables and when there is not the type we expect or the values we want, throw an error
// don't start the server at all
// TEMPMAIL_DOMAINS=domain1.com,domain2.com
// TEMPMAIL_EXPIRATION_ENABLED=true
// TEMPMAIL_EXPIRATION_TIME=48h
// 2. write handlers for dealing with emails with headers and stuff. So that when we are retrieving the messages from the
// mailboxes, we just save the stress of parsing the headers on the fly
// Backend implements SMTP server methods.

type Backend struct {
	db *gorm.DB
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{db: bkd.db}, nil
}

// A Session is returned after EHLO.
type Session struct {
	From string
	To   []string
	db   *gorm.DB
}

func (s *Session) Session() {
	s.From = ""
	s.To = []string{}
}

// We'll implement the Session methods next
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("Mail from:", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, _ *smtp.RcptOptions) error {
	fmt.Println("Rcpt to:", to)
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Logout() error {
	fmt.Println("Logging out")
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := io.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Received message:", string(b))

		msg, err := mail.ReadMessage(bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("failed to parse email: %w", err)
		}

		headersJSON, err := json.Marshal(msg.Header)
		if err != nil {
			return fmt.Errorf("failed to marshal headers: %w", err)
		}

		var address models.Address
		if err := s.db.Where("email = ?", s.To[0]).First(&address).Error; err != nil {
			return fmt.Errorf("Receiver not found")
		}

		message := models.Message{
			Body:        string(b),
			FromAddress: s.From,
			ToAddressID: address.ID,
			ReceivedAt:  time.Now(),
			Subject:     msg.Header.Get("Subject"),
			Headers:     headersJSON,
		}

		if err := s.db.Create(&message).Error; err != nil {
			return err
		}

		return nil
	}
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "testuser" || password != "testpass" {
		return fmt.Errorf("Invalid username or password")
	}

	return nil
}

func (s *Session) Reset() {
	s.From = ""
	s.To = []string{}
}

func (s *Session) Quit() error {
	fmt.Println("Quitting")
	return nil
}

// func init() {
// 	fmt.Println("Initializing mail server and setting up database")
// 	models.SetupDatabase()
// }

func startMailServer() {
	// Create a new backend
	db := models.SetupDatabase()
	backend := &Backend{db: db}

	//Read config from the environment
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env not found, reading system variables")
	}

	requiredVars := []string{"TEMPMAIL_DOMAINS", "TEMPMAIL_EXPIRATION_ENABLED", "TEMPMAIL_EPIRATION_TIME"}
	missingVars := []string{}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %v", missingVars)
		log.Fatalf("Server going down")
		return
	}

	_ = os.Getenv("TEMPMAIL_DOMAINS")

	// Create a new SMTP server
	server := smtp.NewServer(backend)

	// Set the server's address
	server.Addr = ":2525"

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startMailServer()
}

// todo today
// validate email and store in the database
// figuring out the emails will get purged after their duration.
// we look into initial api
//   - generating a random email
//   - viewing the content in the mailbox
// clearing the emails and addresses after they have expired

// headers
//
// email content
// attachments
