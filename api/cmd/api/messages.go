package main

import (
	"fmt"
	"net/http"
	"strconv"

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

func (api *application) getMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := api.readIDParam(r)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	message, err := api.store.Message.Get(idInt)
	if err == store.ErrRecordNotFound {
		api.notFoundResponse(w, r)
		return
	} else if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	err = api.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		api.writeJSONError(w, r)
		return
	}
}

// func (api *application) getMessageByFromAddress(w http.ResponseWriter, r *http.Request) {
// 	email, err := api.readEmailParam(r)
// 	if err != nil {
// 		api.badRequestResponse(w, r, err)
// 		return
// 	}

// 	messages, err := api.store.Message.GetByFromAddress(email)
// 	if err != nil {
// 		api.serverErrorResponse(w, r, err)
// 		return
// 	}

// 	err = api.writeJSON(w, http.StatusOK, envelope{"messages": messages}, nil)
// 	if err != nil {
// 		api.writeJSONError(w, r)
// 		return
// 	}
// }

func (api *application) getMessageByRecepientHandler(w http.ResponseWriter, r *http.Request) {
	email, err := api.readEmailParam(r)
	fmt.Println(email)
	if err != nil {
		api.logger.Println("error here")
		api.badRequestResponse(w, r, err)
		return
	}

	toAddressID, err := api.store.Address.GetByEmail(email)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	messages, err := api.store.Message.GetAll(toAddressID.ID)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	err = api.writeJSON(w, http.StatusOK, envelope{"messages": messages}, nil)
	if err != nil {
		api.writeJSONError(w, r)
		return
	}
}

func (api *application) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := api.readIDParam(r)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		api.badRequestResponse(w, r, err)
		return
	}

	err = api.store.Message.Delete(idInt)
	if err == store.ErrRecordNotFound {
		api.notFoundResponse(w, r)
		return
	} else if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	err = api.writeJSON(w, http.StatusOK, envelope{"message": "message deleted successfully"}, nil)
	if err != nil {
		api.writeJSONError(w, r)
		return
	}
}
