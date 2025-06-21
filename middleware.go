package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/auth"
)

func (app *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prepareResponse(w)
		log.Println(r.URL)

		requestData, err := io.ReadAll(r.Body)
		if err != nil {
			sendError(w, Error{400, "Could not read request body", "Bad Request"}, err)
			return
		}
		r.Body.Close()

		input := struct {
			Token string `json:"token"`
		}{}

		if err := json.Unmarshal(requestData, &input); err != nil {
			sendError(w, Error{400, "Could not read token", "Bad Request"}, err)
			return
		}

		id, err := auth.ValidateSession(app.CACHE, input.Token)
		if err != nil {
			sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
			return
		}

		log.Printf("User authenticated: %d - %s", id, input.Token)

		r.Body = io.NopCloser(bytes.NewReader(requestData))

		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
