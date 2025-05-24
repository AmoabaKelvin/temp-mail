package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-smtp"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
	"github.com/AmoabaKelvin/temp-mail/internal/store"
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
	store *store.Storage
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{store: bkd.store}, nil
}

// A Session is returned after EHLO.
type Session struct {
	From  string
	To    []string
	store *store.Storage
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

// parseMessageBody tries to extract a preferred textual body part (HTML then plain text)
// from the email body. If parsing fails or no suitable part is found,
// it returns the original raw body bytes as a string and its original content type.
func parseMessageBody(originalBodyBytes []byte, header mail.Header) (bodyOutput string, finalContentType string) {
	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Malformed Content-Type ('%s'): %v. Returning raw body.", contentType, err)
		return string(originalBodyBytes), contentType // Return original content type, even if malformed
	}

	var plainTextBody string
	var htmlBody string

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			log.Println("Multipart message lacks boundary. Returning raw body.")
			return string(originalBodyBytes), mediaType
		}

		mr := multipart.NewReader(bytes.NewReader(originalBodyBytes), boundary)
		for {
			p, partErr := mr.NextPart()
			if partErr == io.EOF {
				break
			}
			if partErr != nil {
				log.Printf("Error reading multipart part: %v. Returning raw body.", partErr)
				return string(originalBodyBytes), mediaType
			}
			defer p.Close()

			partMediaTypeStr := p.Header.Get("Content-Type")
			partMediaType, _, partParseErr := mime.ParseMediaType(partMediaTypeStr)
			if partParseErr != nil {
				log.Printf("Skipping part with malformed Content-Type ('%s'): %v", partMediaTypeStr, partParseErr)
				continue
			}

			partBodyBytes, readErr := io.ReadAll(p)
			if readErr != nil {
				log.Printf("Failed to read part body for Content-Type '%s': %v. Skipping part.", partMediaTypeStr, readErr)
				continue
			}

			// A more robust solution would handle charset conversion here
			// e.g., using golang.org/x/net/html/charset based on partParams["charset"]

			if partMediaType == "text/html" {
				htmlBody = string(partBodyBytes)
			} else if partMediaType == "text/plain" {
				plainTextBody = string(partBodyBytes)
			}
		}

		if htmlBody != "" {
			return htmlBody, "text/html"
		}
		if plainTextBody != "" {
			return plainTextBody, "text/plain"
		}

		log.Println("No text/html or text/plain part found in multipart message. Storing raw body.")
		return string(originalBodyBytes), mediaType

	} else if mediaType == "text/plain" || mediaType == "text/html" {
		return string(originalBodyBytes), mediaType
	}

	log.Printf("Content-Type '%s' is not multipart or simple text. Storing raw body.", mediaType)
	return string(originalBodyBytes), mediaType
}

func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read incoming data: %w", err)
	}
	log.Printf("Received raw message, %d bytes", len(b))

	msg, err := mail.ReadMessage(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to parse email: %w", err)
	}

	headersJSON, err := json.Marshal(msg.Header)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	ctx := context.Background()
	address, err := s.store.Addresses.Get(ctx, s.To[0])
	if err != nil {
		return fmt.Errorf("receiver address '%s' not found: %w", s.To[0], err)
	}

	if address.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("receiver address '%s' has expired", s.To[0])
	}

	rawMessageBodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to read raw message body: %w", err)
	}

	finalBodyString, _ := parseMessageBody(rawMessageBodyBytes, msg.Header)

	log.Printf("Storing message for %s, Subject: %s, Body length: %d",
		s.To[0], msg.Header.Get("Subject"), len(finalBodyString))

	message := store.Message{
		Body:        finalBodyString,
		FromAddress: s.From,
		ToAddressID: uint(address.ID),
		ReceivedAt:  time.Now(),
		Subject:     msg.Header.Get("Subject"),
		Headers:     headersJSON,
	}

	if err := s.store.Messages.Create(ctx, &message); err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	log.Printf("Successfully stored message ID %d for %s", message.ID, s.To[0])
	return nil
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "testuser" || password != "testpass" {
		return fmt.Errorf("invalid username or password")
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

func startMailServer() {
	db, err := db.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	store := store.NewStorage(db)

	backend := &Backend{store: store}

	server := smtp.NewServer(backend)

	server.Addr = "0.0.0.0:25"

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startMailServer()
}
