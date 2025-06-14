package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"

	"server/database"
)

func GenerateSecureToken(length int64) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

func CreateSession(db *sql.DB, id int64) (string, error) {
	token, err := GenerateSecureToken(64)
	if err != nil {
		return "", err
	}

	if err := database.InsertToken(db, id, token); err != nil {
		return "", nil
	}

	return token, nil
}

func ValidateSession(db *sql.DB, token string) (int64, error) {
	id, err := database.GetToken(db, token)

	return id, err
}
