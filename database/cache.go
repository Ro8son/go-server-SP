package database

import (
	"database/sql"
	"errors"
	"time"
)

func InsertToken(db *sql.DB, login, token string) error {
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := db.Exec(`
        INSERT INTO sessions (login, token, expires_at)
        VALUES (?, ?, ?)`,
		login, token, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func GetToken(db *sql.DB, token string) (string, time.Time, error) {
	var login string
	var expiresAt time.Time

	err := db.QueryRow(`
        SELECT login, expires_at FROM sessions 
        WHERE token = ?`, token).Scan(&login, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", time.Time{}, errors.New("invalid session")
		}
		return "", time.Time{}, err
	}

	return login, expiresAt, err
}

func DeleteToken(db *sql.DB, token string) {
	_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
