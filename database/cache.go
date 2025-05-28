package database

import (
	"database/sql"
	"errors"
)

func SetupCache(db *sql.DB) error {
	// Create sessions table
	_, err := db.Exec(`
    CREATE TABLE sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        login TEXT NOT NULL,
        token TEXT NOT NULL UNIQUE
    )`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
    CREATE TABLE uploads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT NOT NULL,
				transaction_id INTEGER NOT NULL UNIQUE
    		file_name TEXT NOT NULL,
    		title TEXT,
    		description TEXT,
    		coordinates TEXT
    )`)

	return nil
}

func InsertUploadMeta(db *sql.DB, id, token, fileName, title, description, coordinates string) (int, error) {
	//expiresAt := time.Now().Add(24 * time.Hour)

	result, err := db.Exec(`
        INSERT INTO uploads (token, transaction_id, file_name, title, description, coordinates)
        VALUES (?, ?, ?, ?, ?, ?)`, token, id, fileName, title, description, coordinates)
	if err != nil {
		return -1, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(lastInsertId), err
}

func GetMetadataId(db *sql.DB, transaction_id string, token string) (int, error) {
	var id int

	err := db.QueryRow(`
        SELECT id FROM uploads 
        WHERE token = ? AND transaction_id = ?`, token, transaction_id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("invalid transaction id or token")
		}
		return -1, err
	}

	return id, err
}

func InsertToken(db *sql.DB, login, token string) error {

	_, err := db.Exec(`
        INSERT INTO sessions (login, token)
        VALUES (?, ?)`,
		login, token)
	if err != nil {
		return err
	}

	return nil
}

func GetToken(db *sql.DB, token string) (string, error) {
	var login string

	err := db.QueryRow(`
        SELECT login FROM sessions 
        WHERE token = ?`, token).Scan(&login)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid session")
		}
		return "", err
	}

	return login, err
}

func DeleteToken(db *sql.DB, token string) {
	_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
