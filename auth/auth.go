package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"server/database"
	"time"
)

func generateSecureToken(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

func CreateSession(db *sql.DB, login string) (string, error) {
	token, err := generateSecureToken(64)
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

	login, expiresAt, err := database.GetToken(db, token)
	if err != nil {
		log.Println(err)
	}

	if time.Now().After(expiresAt) {
		database.DeleteToken(db, token)
		return "", errors.New("session expired")
	}

	return login, nil
}
