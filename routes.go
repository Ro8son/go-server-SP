package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.register)
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)
	router.HandleFunc("POST /file/upload", app.uploadFile)
	router.HandleFunc("GET /file/list", app.getFileList)
	router.HandleFunc("GET /file/download", app.fileDownload)

	return router
}
