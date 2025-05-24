package main

import "net/http"

func (app *application) notFound(w http.ResponseWriter) {
	app.writeErrorJSON(w, http.StatusNotFound, "not found")
}

func (app *application) serverError(w http.ResponseWriter) {
	app.writeErrorJSON(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequest(w http.ResponseWriter, message string) {
	app.writeErrorJSON(w, http.StatusBadRequest, message)
}
