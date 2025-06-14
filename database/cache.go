package database

import (
	"database/sql"
	"errors"
	"log"
)

func SetupCache(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE sessions ( id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        token TEXT NOT NULL UNIQUE
    )`)

	return err
}

func InsertToken(db *sql.DB, id int64, token string) error {
	query := `INSERT INTO sessions (user_id, token) VALUES (?, ?)`

	_, err := db.Exec(query, id, token)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetToken(db *sql.DB, token string) (int64, error) {
	var id int64

	err := db.QueryRow(`
        SELECT user_id FROM sessions
        WHERE token = ?`, token).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("invalid session")
		}
		return -1, err
	}

	return id, err
}

func DeleteToken(db *sql.DB, token string) {
	_, _ = db.Exec("DELETE FROM sessions WHERE token = ?", token)
}
