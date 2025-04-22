package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
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

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(`
        INSERT INTO sessions (login, token, expires_at)
        VALUES (?, ?, ?)`,
		login, token, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateSession(db *sql.DB, token string) (string, error) {
	var login string
	var expiresAt time.Time

	log.Println(token)

	err := db.QueryRow(`
        SELECT login, expires_at FROM sessions 
        WHERE token = ?`, token).Scan(&login, &expiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid session")
		}
		return "", err
	}

	if time.Now().After(expiresAt) {
		_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
		return "", errors.New("session expired")
	}

	return login, nil
}
