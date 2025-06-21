package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /register", app.register)
	router.Handle("PUT /register", app.authenticate(http.HandlerFunc(app.updateUser)))
	router.HandleFunc("POST /login", app.login)
	router.HandleFunc("POST /logout", app.logout)

	router.Handle("POST /file/upload", app.authenticate(http.HandlerFunc(app.uploadFile)))
	router.Handle("POST /file/share/add", app.authenticate(http.HandlerFunc(app.shareFile)))
	router.Handle("POST /file/share/get", app.authenticate(http.HandlerFunc(app.getShareFile)))
	router.Handle("POST /file/download", app.authenticate(http.HandlerFunc(app.fileDownload)))
	router.Handle("POST /file/list", app.authenticate(http.HandlerFunc(app.getFileList)))
	router.Handle("POST /file/tags", app.authenticate(http.HandlerFunc(app.getTags)))
	router.Handle("GET /shared/{id}/{pass}", http.HandlerFunc(app.downloadSharedFile))

	router.Handle("POST /album/add", app.authenticate(http.HandlerFunc(app.addAlbum)))
	router.Handle("POST /album/list", app.authenticate(http.HandlerFunc(app.getAlbums)))
	router.Handle("POST /album/addFile", app.authenticate(http.HandlerFunc(app.addFileToAlbum)))
	router.Handle("POST /album/getFile", app.authenticate(http.HandlerFunc(app.addFileToAlbum)))

	return router
}
