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
	usr "server/user"

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

	found, err := database.GetUser(app.DB, user.Login)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
		return
	}
	if found != "" {
		sendError(w, Error{418, "No tea for this User", "I'm a teapot"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"})
		return
	}

	user.Login, err = usr.AddUser(app.DB, user.Login, string(hashedPassword))
	if err != nil {
		log.Println(err)
		sendError(w, Error{500, "Could not add user", "Internal Server Error"})
		return
	}

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
		return
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

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	token := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	database.DeleteToken(app.CACHE, token.Token)

	token.Token = "" // should be changed
	if err := json.NewEncoder(w).Encode(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *app) initFileUpload(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	metadata := struct {
		Token    string `json:"token"`
		Login    string `json:"login"`
		FileName string `json:"file_name"`
		Id       string `json:"transaction_id"`
		// some other data (soon™)
	}{}

	err := json.NewDecoder(r.Body).Decode(&metadata)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	metadata.Login, err = auth.ValidateSession(app.CACHE, metadata.Token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
		return
	}

	metadata.Id, err = auth.GenerateSecureToken(128)

	err = database.InsertUploadMeta(app.CACHE, metadata.Id, metadata.Token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Database", "Internal Server Error"})
		return
	}

	metadata.Token = ""
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (app *app) fileUpload(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")
	id := r.FormValue("transaction_id")

	login, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
		return
	}

	_, err = database.GetUploadMetadata(app.CACHE, id, token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "No metdata found", "Bad Request"})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire the file", "Bad Request"})
		return
	}
	defer file.Close()

	f, err := os.OpenFile("../storage/users/"+login+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire file path", "Internal Server Error"})
		return
	}

	io.Copy(f, file)

	f.Close()
}

func (app *app) getFileList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	files := struct {
		Token string   `json:"token"`
		Files []string `json:"files"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&files); err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"})
		return
	}

	login, err := auth.ValidateSession(app.CACHE, files.Token)
	if err != nil {
		log.Println(err)
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"})
		return
	}

	entries, err := os.ReadDir("../storage/users/" + login)
	if err != nil {
		log.Println(err)
		sendError(w, Error{400, "Could not acquire file path", "Internal Server Error"})
		return
	}

	for _, e := range entries {
		files.Files = append(files.Files, e.Name())
	}

	files.Token = ""
	if err = json.NewEncoder(w).Encode(files); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
