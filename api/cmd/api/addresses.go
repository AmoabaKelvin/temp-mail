package main

import (
	"net/http"

	"github.com/AmoabaKelvin/temp-mail/api/internal/store"
	"github.com/AmoabaKelvin/temp-mail/api/internal/validator"
)

func (api *application) generateAddressHandler(w http.ResponseWriter, r *http.Request) {
	addressSchema := api.generateAddress()

	address := &store.Address{
		Email:     addressSchema.Email,
		ExpiresAt: addressSchema.ExpiresAt,
	}
	v := validator.NewValidator()

	if store.ValidateAddress(v, address); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}

	err := api.store.Address.Insert(address)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}

	err = api.writeJSON(w, http.StatusCreated, envelope{"address": address}, nil)
	if err != nil {
		api.writeJSONError(w, r)
		return
	}
}
