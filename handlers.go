package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"server/auth"
	"server/database"

	_ "github.com/glebarez/go-sqlite"
	"golang.org/x/crypto/bcrypt"
)

func prepareResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
}

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorType  string `json:"error_type"`
	// to be expanded
}

// interface for json needed
func sendError(w http.ResponseWriter, error Error) {
	if err := json.NewEncoder(w).Encode(error); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	user := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"})
		return
	}

	err = database.AddUser(app.DB, user.Login, string(hashedPassword))
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not add user", "Internal Server Error"})
		return
	}
	log.Printf("Added User: \nLogin: %s\nPassword: %s", user.Login, hashedPassword)

	user.Password = strings.Repeat("*", len(user.Password)) // should be changed
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	user := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	hashedPassword, err := database.GetUser(app.DB, user.Login)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		log.Println(err)
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"})
		return
	} else {
		token, err := auth.CreateSession(app.CACHE, user.Login)
		if err != nil {
			log.Println(err)
			sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"})
			return
		}

		log.Printf("User: %s - Logged in with token: %s", user.Login, token)

		user.Token = token
		user.Password = strings.Repeat("*", len(user.Password)) // should be changed
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (app *app) FileUpload(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)

	token := r.FormValue("token")

	login, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		w.Write([]byte("Invalid token"))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	if err := os.MkdirAll("./"+login, os.ModePerm); err != nil {
		log.Println(err)
	}

	f, err := os.OpenFile("./"+login+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
	}

	io.Copy(f, file)

	w.Write([]byte("file uploaded as user: " + login))
}

func (app *app) getFileList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)

	token := r.FormValue("token")

	login, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		w.Write([]byte("Invalid token"))
		return
	}

	entries, err := os.ReadDir("./" + login)
	if err != nil {
		log.Println(err)
	}

	var output string

	for _, e := range entries {
		output += e.Name() + " "
	}

	w.Write([]byte(output))
}
