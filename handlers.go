package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"

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
func sendError(w http.ResponseWriter, error Error, err error) {
	log.Println(err)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (app *app) register(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	_, found, _, err := database.GetUser(app.DB, input.Login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	} else if found != "" {
		sendError(w, Error{418, "No tea for this User", "I'm a teapot"}, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"}, err)
		return
	}

	input.Login, err = usr.AddUser(app.DB, input.Login, string(hashedPassword), input.Email)
	if err != nil {
		sendError(w, Error{500, "Could not add user", "Internal Server Error"}, err)
		return
	}

	log.Printf("Add user: -- Login: %s - Password: %s", input.Login, input.Password)

	w.WriteHeader(http.StatusOK)
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	output := struct {
		Token   string `json:"token"`
		IsAdmin int    `json:"is_admin"`
	}{}

	var hashedPassword string

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	_, hashedPassword, output.IsAdmin, err = database.GetUser(app.DB, input.Login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"}, err)
		return
	} else {
		output.Token, err = auth.CreateSession(app.CACHE, input.Login)
		if err != nil {
			sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"}, err)
			return
		}

		log.Printf("Login -- Login: %s - Token: %s", input.Login, output.Token)

		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	input := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	database.DeleteToken(app.CACHE, input.Token)

	log.Printf("Logout: -- Removed: %s (maybe valid)", input.Token)

	w.WriteHeader(http.StatusOK)
}

func (app *app) uploadFile(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	type file struct {
		File        string `json:"file"`
		FileName    string `json:"file_name"`
		Title       string `json:"title"`       //optional
		Description string `json:"description"` //optional
		Coordinates string `json:"coordinates"` //optional
	}

	input := struct {
		Token string `json:"token"`
		Files []file `json:"files"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	login, err := auth.ValidateSession(app.CACHE, input.Token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	id, _, _, err := database.GetUser(app.DB, login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for x := range input.Files {
		fileId, err := database.AddFile(app.DB, id, input.Files[x].FileName, input.Files[x].Title, input.Files[x].Description, input.Files[x].Coordinates)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		data, err := base64.StdEncoding.DecodeString(input.Files[x].File)
		if err != nil {
			sendError(w, Error{400, "Decoding", "Internal Server Error"}, err)
			return
		}

		fileIdStr := strconv.FormatInt(fileId, 16)

		f, err := os.OpenFile("../storage/users/"+login+"/"+fileIdStr, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			sendError(w, Error{400, "Could not acquire file path", "Internal Server Error"}, err)
			return
		}

		f.Write(data)
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getFileList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	output := struct {
		Files []database.File `json:"files"`
	}{}

	r.ParseMultipartForm(32 << 20)
	token := r.FormValue("token")

	login, err := auth.ValidateSession(app.CACHE, token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	id, _, _, err := database.GetUser(app.DB, login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Files, err = database.GetFileTitles(app.DB, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = json.NewEncoder(w).Encode(output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *app) fileDownload(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token   string  `json:"token"`
		FileIds []int64 `json:"file_ids"`
	}{}

	output := struct {
		Files []database.File `json:"files"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	login, err := auth.ValidateSession(app.CACHE, input.Token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	id, _, _, err := database.GetUser(app.DB, login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Files, err = database.GetFileTitles(app.DB, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for i := range output.Files {
		if slices.Contains(input.FileIds, output.Files[i].Id) {
			file, err := os.ReadFile("../storage/users/" + login + "/" + strconv.FormatInt(output.Files[i].Id, 16))
			if err != nil {
				sendError(w, Error{400, "Error opening file:" + output.Files[i].FileName, "Internal Server Error"}, err)
				return
			}

			output.Files[i].File = base64.StdEncoding.EncodeToString(file)
		}
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
