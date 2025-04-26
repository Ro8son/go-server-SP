package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.register)
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)
	router.HandleFunc("POST /upload", app.fileUpload)
	router.HandleFunc("GET /upload", app.getFileList)
	router.HandleFunc("POST /metadata", app.initFileUpload)

	return router
}
