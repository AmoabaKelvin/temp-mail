package main

import (
	"encoding/json"
	"maps"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	w.Header().Set("Content-Type", "application/json")
	maps.Copy(w.Header(), headers)

	w.WriteHeader(status)

	type Envelope struct {
		Data any `json:"data"`
	}

	envelope := Envelope{
		Data: data,
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(envelope)
}

// func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
// 	maxBytes := 1_048_576
// 	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	return decoder.Decode(data)
// }

func (app *application) writeErrorJSON(w http.ResponseWriter, status int, message any) error {
	return app.writeJSON(w, status, map[string]any{"error": message}, nil)
}
