package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AmoabaKelvin/temp-mail/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getMessages(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if email == "" {
		app.badRequest(w, "email parameter is required")
		return
	}

	address, err := app.store.Addresses.Get(r.Context(), email)
	if errors.Is(err, store.ErrNotFound) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w)
		return
	}

	messages, err := app.store.Messages.Get(r.Context(), address.ID)
	if err != nil {
		app.serverError(w)
		return
	}

	app.writeJSON(w, http.StatusOK, messages, nil)
}

func (app *application) deleteMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequest(w, "invalid message ID")
		return
	}

	if err := app.store.Messages.Delete(r.Context(), id); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFound(w)
		default:
			app.serverError(w)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Message deleted successfully"}, nil)
}

func (app *application) updateMessageReadAt(w http.ResponseWriter, r *http.Request) {
	readAt := time.Now()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequest(w, "invalid message ID")
		return
	}

	err = app.store.Messages.SetReadAt(r.Context(), id, &readAt)
	if err != nil {
		app.serverError(w)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"message": "Message read status updated successfully"}, nil)
}
