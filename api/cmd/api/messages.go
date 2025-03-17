package main

import (
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/api/internal/store"
	"github.com/AmoabaKelvin/temp-mail/api/internal/validator"
)

func (api *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input = store.MessageSchema{}

	err := api.readJSON(w, r, &input)
	if err != nil {
		api.logger.Println(err)
		api.badRequestResponse(w, r, err)
		return
	}

	message := &store.Message{
		FromAddress: input.FromAddress,
		ToAddressID: input.ToAddressID,
		Subject:     input.Subject,
		Body:        input.Body,
		ReceivedAt:  input.ReceivedAt,
	}

	v := validator.NewValidator()

	if store.ValidateMessage(v, message); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = api.store.Message.Insert(message)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	err = api.writeJSON(w, http.StatusCreated, envelope{"message": message}, nil)
	if err != nil {
		api.writeJSONError(w, r)
		return
	}
}
