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

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.register)
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)

	router.Handle("POST /file/upload", app.authenticate(http.HandlerFunc(app.uploadFile)))
	router.Handle("POST /file/share/add", app.authenticate(http.HandlerFunc(app.shareFile)))
	router.Handle("POST /file/share/get", app.authenticate(http.HandlerFunc(app.getShareFile)))
	router.Handle("POST /file/download", app.authenticate(http.HandlerFunc(app.fileDownload)))
	router.Handle("POST /file/list", app.authenticate(http.HandlerFunc(app.getFileList)))
	router.Handle("GET /shared/{id}/{pass}", http.HandlerFunc(app.downloadSharedFile))

	router.Handle("POST /album/add", app.authenticate(http.HandlerFunc(app.addAlbum)))
	router.HandleFunc("POST /album/list", app.fileDownload)
	router.HandleFunc("POST /album/del", app.fileDownload)

	return router
}

func (app *app) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
