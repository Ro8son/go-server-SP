package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.register)
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /upload", app.FileUpload)
	router.HandleFunc("GET /upload", app.getFileList)

	return router
}
