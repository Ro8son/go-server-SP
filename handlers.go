package main

import (
	"database/sql"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		sendError(w, Error{500, "Could not generate hash from password", "Internal Server Error"}, err)
		return
	}

	input.Login, err = usr.AddUser(app.Query, input.Login, string(hashedPassword), input.Email)
	if err != nil {
		sendError(w, Error{500, "Could not add user", "Internal Server Error"}, err)
		return
	}

	log.Printf("Add user: -- Login: %s - Password: %s", input.Login, input.Password)

	w.WriteHeader(http.StatusOK)
}

func (app *app) updateUser(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	var input database.UpdateUserParams
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	output, err := app.Query.UpdateUser(app.Ctx, input)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output.Password = ""
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	output := struct {
		Token   string `json:"token"`
		Profile string `json:"profile"`
		Email   string `json:"email"`
		IsAdmin int    `json:"is_admin"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	user, err := app.Query.GetUserByLogin(app.Ctx, input.Login)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		sendError(w, Error{401, "Wrong password or login", "Unauthorized"}, err)
		return
	} else {
		output.Token, err = auth.CreateSession(app.CACHE, user.ID)
		if err != nil {
			sendError(w, Error{500, "Could not generate a new token", "Internal Server Error"}, err)
			return
		}

		profile, err := app.Query.GetProfile(app.Ctx, user.ID)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		email, err := app.Query.GetEmail(app.Ctx, user.ID)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		log.Printf("Login -- Login: %s - Token: %s", input.Login, output.Token)

		output.IsAdmin = int(user.IsAdmin)
		output.Profile = profile.String
		output.Email = email.String

		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (app *app) logout(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
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
		File     string                 `json:"file"`
		Metadata database.AddFileParams `json:"metadata"`
		Tags     []string               `json:"tags"`
	}

	input := struct {
		Files []file `json:"files"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	id := r.Context().Value("id").(int64)

	user, err := app.Query.GetUser(app.Ctx, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for _, file := range input.Files {
		file.Metadata.OwnerID = id
		id, err := app.Query.AddFile(app.Ctx, file.Metadata)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		data, err := base64.StdEncoding.DecodeString(file.File)
		if err != nil {
			sendError(w, Error{400, "Decoding", "Internal Server Error"}, err)
			return
		}

		fileIdStr := strconv.FormatInt(id, 16)

		f, err := os.OpenFile("../storage/users/"+user.Login+"/"+fileIdStr, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			sendError(w, Error{400, "Could not acquire file path", "Internal Server Error"}, err)
			return
		}

		f.Write(data)

		for _, tag := range file.Tags {
			tagDB, err := app.Query.GetTagByName(app.Ctx, tag)
			if err != nil && err != sql.ErrNoRows {
				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
				return
			}

			if sql.ErrNoRows == err {
				tagDB.ID, err = app.Query.AddTag(app.Ctx, tag)
			}

			err = app.Query.TagsConnect(app.Ctx, database.TagsConnectParams{FileID: id, TagID: tagDB.ID})
			if err != nil && err != sql.ErrNoRows {
				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
				return
			}
		}

	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getFileList(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	id := r.Context().Value("id").(int64)
	type File struct {
		File database.GetFilesRow `json:"file"`
		Tags []string             `json:"tags"`
	}

	output := struct {
		File []File `json:"file"`
	}{}

	user, err := app.Query.GetUser(app.Ctx, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	files, err := app.Query.GetFiles(app.Ctx, user.ID)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for _, file := range files {
		tags, err := app.Query.GetTagsByFile(app.Ctx, file.ID)
		if err != nil {
			sendError(w, Error{400, "Database", "Internal Server Error"}, err)
			return
		}

		var tagNames []string
		for _, tagID := range tags {
			tagName, err := app.Query.GetTagById(app.Ctx, tagID)
			if err != nil {
				sendError(w, Error{400, "Database", "Internal Server Error"}, err)
				return
			}
			tagNames = append(tagNames, tagName.Name)
		}
		output.File = append(output.File, File{File: file, Tags: tagNames})
	}

	if err = json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type File struct {
	Id       int64  `json:"id"`
	FileName string `json:"file_name"`
	File     string `json:"file"`
	// checksum
}

func (app *app) fileDownload(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)

	input := struct {
		Token   string  `json:"token"`
		FileIds []int64 `json:"file_ids"`
	}{}

	output := struct {
		Files []File `json:"files"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	id, err := auth.ValidateSession(app.CACHE, input.Token)
	if err != nil {
		sendError(w, Error{401, "Incorrect Token", "Unauthorized"}, err)
		return
	}

	user, err := app.Query.GetUser(app.Ctx, int64(id))
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	files, err := app.Query.GetFiles(app.Ctx, user.ID)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	for i := range files {
		if slices.Contains(input.FileIds, files[i].ID) {
			file, err := os.ReadFile("../storage/users/" + user.Login + "/" + strconv.FormatInt(files[i].ID, 16))
			if err != nil {
				sendError(w, Error{400, "Error opening file:" + output.Files[i].FileName, "Internal Server Error"}, err)
				return
			}

			data := base64.StdEncoding.EncodeToString(file)
			output.Files = append(output.Files, File{Id: files[i].ID, FileName: files[i].FileName, File: data})
		}
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getTags(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	output, err := app.Query.GetTags(app.Ctx)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addAlbum(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	intput := struct {
		AlbumTitle database.AddAlbumParams `json:"album_title"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&intput); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}
	intput.AlbumTitle.OwnerID = r.Context().Value("id").(int64)

	if err := app.Query.AddAlbum(app.Ctx, intput.AlbumTitle); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getAlbums(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	id := r.Context().Value("id").(int64)

	output, err := app.Query.GetAlbums(app.Ctx, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) addFileToAlbum(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	var input database.AddToAlbumParams
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	if err := app.Query.AddToAlbum(app.Ctx, input); err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *app) getFileFromAlbum(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	input := struct {
		AlbumID int64 `json:"album_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	output, err := app.Query.GetFileFromAlbum(app.Ctx, input.AlbumID)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) shareFile(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	var input database.AddGuestFileParams

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, Error{400, "Could not acquire json data", "Bad Request"}, err)
		return
	}

	url, err := auth.GenerateSecureToken(32)
	if err != nil {
		sendError(w, Error{400, "Could not generate URL", "Internal Server Error"}, err)
		return
	}
	input.Url = url

	share, err := app.Query.AddGuestFile(app.Ctx, input)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	output := struct {
		Url string `json:"url"`
	}{
		Url: "shared/" + strconv.Itoa(int(share.ID)) + "/" + share.Url,
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) getShareFile(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	id := r.Context().Value("id").(int64)

	output, err := app.Query.GetSharedFiles(app.Ctx, id)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&output); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *app) downloadSharedFile(w http.ResponseWriter, r *http.Request) {
	prepareResponse(w)
	id, err := strconv.Atoi(r.PathValue("id"))
	pass := r.PathValue("pass")

	var input database.GetShareDownloadParams
	input.ID = int64(id)
	input.Url = pass

	output, err := app.Query.GetShareDownload(app.Ctx, input)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	login, err := app.Query.GetLogin(app.Ctx, output.OwnerID.Int64)
	if err != nil {
		sendError(w, Error{400, "Database", "Internal Server Error"}, err)
		return
	}

	file, err := os.ReadFile("../storage/users/" + login + "/" + strconv.FormatInt(output.ID.Int64, 16))
	if err != nil {
		sendError(w, Error{400, "Error opening file:" + output.FileName.String, "Internal Server Error"}, err)
		return
	}

	w.Write(file)
}
