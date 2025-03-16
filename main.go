package main

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/AmoabaKelvin/temp-mail/pkg/dto"
)

// address generation service
// smtp receiver: receive incoming mails, parse and store them in a database
// db layer: interact with the database to store and retrieve data
// api layer: expose the service through an API
// generated addresses need to expire after a certain time
// cleanup layer: clean up expired addresses from the database
// webhooks
// kelvinamoaba@something.com : MX record -> can point to an IP address

// https://github.com/matoous/go-nanoid
func generateAddress() models.Address {
	id, err := gonanoid.New()
	domain := "example.com" // todo: change this later to the actual domain
	if err != nil {
		panic(err)
	}
	return models.Address{
		Email:     fmt.Sprintf("%s@%s", id, domain),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
}

// ui -> generate -> generateAddress() -> db

func main() {
	fmt.Println(generateAddress())
	// check if there is a collision (skip)
	// store that information somewhere (db, temp: struct)
	// return the address back to the caller
	// qVOv2Id_eOszwvERy6bok@example.com
	//

	// address := models.Address{
	// 	ID:        generateAddress(),
	// 	ExpiresAt: time.Now().Add(time.Hour * 24),
}

// define the db models.
// add in helper functions to interact with the database (shared)
// parsing the received emails
// storing them
