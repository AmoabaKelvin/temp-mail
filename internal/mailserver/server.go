package mailserver

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
	"strings"
	"time"

	"github.com/emersion/go-smtp"

	"github.com/AmoabaKelvin/temp-mail/internal/db"
	"github.com/AmoabaKelvin/temp-mail/internal/store"
)

// Backend implements SMTP server methods.
type Backend struct {
	store *store.Storage
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{store: bkd.store}, nil
}

// Session is returned after EHLO.
type Session struct {
	From  string
	To    []string
	store *store.Storage
}

func (s *Session) Session() {
	s.From = ""
	s.To = []string{}
}

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

// readAndParseMessage reads the incoming data and parses it as an email message
func readAndParseMessage(r io.Reader) (*mail.Message, []byte, error) {
	rawData, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read incoming data: %w", err)
	}
	log.Printf("Received raw message, %d bytes", len(rawData))

	msg, err := mail.ReadMessage(bytes.NewReader(rawData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse email: %w", err)
	}

	return msg, rawData, nil
}

// validateRecipient checks if the recipient address exists and is not expired
func validateRecipient(store *store.Storage, address string) (*store.Address, error) {
	ctx := context.Background()
	addr, err := store.Addresses.Get(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("receiver address '%s' not found: %w", address, err)
	}

	if addr.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("receiver address '%s' has expired", address)
	}

	return addr, nil
}

// extractMessageBody reads and processes the message body
func extractMessageBody(msg *mail.Message) (string, error) {
	rawBodyBytes, err := io.ReadAll(msg.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read raw message body: %w", err)
	}

	processedBody, _ := parseMessageBody(rawBodyBytes, msg.Header)
	return processedBody, nil
}

// createMessage constructs a store.Message from the parsed email data
func createMessage(from, subject, body string, headers []byte, addressID uint) store.Message {
	return store.Message{
		Body:        body,
		FromAddress: from,
		ToAddressID: addressID,
		ReceivedAt:  time.Now(),
		Subject:     subject,
		Headers:     headers,
	}
}

// storeMessage persists the message to the database
func storeMessage(storage *store.Storage, message *store.Message) error {
	ctx := context.Background()
	if err := storage.Messages.Create(ctx, message); err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}
	return nil
}

// parseMultipartBody extracts text parts from a multipart message
func parseMultipartBody(bodyBytes []byte, boundary string) (htmlBody, plainBody string, err error) {
	mr := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)

	for {
		part, partErr := mr.NextPart()
		if partErr == io.EOF {
			break
		}
		if partErr != nil {
			return "", "", fmt.Errorf("error reading multipart part: %w", partErr)
		}
		defer part.Close()

		partContentType := part.Header.Get("Content-Type")
		partMediaType, _, parseErr := mime.ParseMediaType(partContentType)
		if parseErr != nil {
			log.Printf("Skipping part with malformed Content-Type ('%s'): %v", partContentType, parseErr)
			continue
		}

		partBodyBytes, readErr := io.ReadAll(part)
		if readErr != nil {
			log.Printf("Failed to read part body for Content-Type '%s': %v. Skipping part.", partContentType, readErr)
			continue
		}

		switch partMediaType {
		case "text/html":
			htmlBody = string(partBodyBytes)
		case "text/plain":
			plainBody = string(partBodyBytes)
		}
	}

	return htmlBody, plainBody, nil
}

// selectPreferredBody chooses the best available body content (HTML preferred over plain text)
func selectPreferredBody(htmlBody, plainBody string) (body, contentType string) {
	if htmlBody != "" {
		return htmlBody, "text/html"
	}
	if plainBody != "" {
		return plainBody, "text/plain"
	}
	return "", ""
}

// parseMessageBody tries to extract a preferred textual body part (HTML then plain text)
// from the email body. If parsing fails or no suitable part is found,
// it returns the original raw body bytes as a string and its original content type.
func parseMessageBody(originalBodyBytes []byte, header mail.Header) (bodyOutput string, finalContentType string) {
	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Malformed Content-Type ('%s'): %v. Returning raw body.", contentType, err)
		return string(originalBodyBytes), contentType
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			log.Println("Multipart message lacks boundary. Returning raw body.")
			return string(originalBodyBytes), mediaType
		}

		htmlBody, plainBody, err := parseMultipartBody(originalBodyBytes, boundary)
		if err != nil {
			log.Printf("Error parsing multipart body: %v. Returning raw body.", err)
			return string(originalBodyBytes), mediaType
		}

		if body, ctype := selectPreferredBody(htmlBody, plainBody); body != "" {
			return body, ctype
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
	// Read and parse the incoming email
	msg, _, err := readAndParseMessage(r)
	if err != nil {
		return err
	}

	// Marshal headers for storage
	headersJSON, err := json.Marshal(msg.Header)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	// Validate recipient address
	address, err := validateRecipient(s.store, s.To[0])
	if err != nil {
		return err
	}

	// Extract and process message body
	body, err := extractMessageBody(msg)
	if err != nil {
		return err
	}

	// Create message object
	subject := msg.Header.Get("Subject")
	message := createMessage(s.From, subject, body, headersJSON, uint(address.ID))

	// Log the operation
	log.Printf("Storing message for %s, Subject: %s, Body length: %d",
		s.To[0], subject, len(body))

	// Store the message
	if err := storeMessage(s.store, &message); err != nil {
		return err
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

// Start initializes and starts the SMTP mail server
func Start(databaseUrl, port string) error {
	database, err := db.New(databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	storage := store.NewStorage(database)

	backend := &Backend{store: storage}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf("0.0.0.0:%s", port)

	log.Printf("Starting SMTP server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start SMTP server: %w", err)
	}

	return nil
}
