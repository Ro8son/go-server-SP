package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"

	"server/database"
)

func GenerateSecureToken(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

func CreateSession(db *sql.DB, login string) (string, error) {
	token, err := GenerateSecureToken(64)
	if err != nil {
		return "", err
	}

	if err = database.InsertToken(db, login, token); err != nil {
		return "", nil
	}

	return token, nil
}

func ValidateSession(db *sql.DB, token string) (string, error) {

	log.Println(token)

	login, err := database.GetToken(db, token)
	if err != nil {
		log.Println(err)
	}

	return login, nil
}
