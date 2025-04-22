package database

import (
	"database/sql"
	"log"
)

func SetupDatabase(db *sql.DB) error {
	// Create sessions table
	_, err := db.Exec(`
    CREATE TABLE sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        login TEXT NOT NULL UNIQUE,
        token TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL
    )`)
	if err != nil {
		return err
	}

	return nil
}

func AddUser(db *sql.DB, login, password string) error {
	query := "INSERT INTO Users VALUES(?, ?, ?)"

	id := getUserCount(db) + 1

	_, err := db.Exec(query, id, login, password)
	return err
}

func getUserCount(db *sql.DB) int {
	var count int
	query := "SELECT COUNT(*) FROM Users"

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println("Error getting user count:", err)
		return 0
	}

	return count
}

func GetUser(db *sql.DB, login string) (string, error) {
	query := "SELECT login, password FROM Users WHERE login = ?"

	var foundLogin, foundPassword string
	err := db.QueryRow(query, login).Scan(&foundLogin, &foundPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", nil
		}
		log.Println("Error retrieving user:", err)
		return "", err
	}

	return foundPassword, nil
}
