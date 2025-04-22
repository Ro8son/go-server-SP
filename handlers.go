package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"server/auth"
	"server/database"

	_ "github.com/glebarez/go-sqlite"
	"golang.org/x/crypto/bcrypt"
)

func prepareResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func getCookie(w http.ResponseWriter, r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.Write([]byte("Cookie not found"))
		} else {
			log.Println("Could not acquire cookie", err)
		}
		return nil, err
	}
	return cookie, nil
}

func setCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	expire := time.Now().Add(duration)
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Expires: expire})
}

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Error parsing post"))
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println(err)
	}

	database.AddUser(app.DB, login, string(hashedPassword))
	log.Printf("Added User: \nLogin: %s\nPassword: %s", login, hashedPassword)

	response := struct {
		Login string `json:"login"`
	}{
		Login: login,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	login := r.FormValue("login")
	password := r.FormValue("password")

	hashedPassword, err := database.GetUser(app.DB, login)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
	} else {
		token, err := auth.CreateSession(app.CACHE, login)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("User: %s - Logged in with token: %s", login, token)

		response := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
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
