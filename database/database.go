package database

import (
	"database/sql"
	"log"
)

// func SetupDatabase(db *sql.DB) error {
// 	// Create sessions table
// 	_, err := db.Exec(`
//     CREATE TABLE Users (
//         id INTEGER PRIMARY KEY AUTOINCREMENT,
//         login TEXT NOT NULL UNIQUE,
//         password TEXT NOT NULL
//     )`)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func AddUser(db *sql.DB, login, password, email string) error {
	query := "INSERT INTO Users (login, password, email) VALUES(?, ?, ?)"

	_, err := db.Exec(query, login, password, email)
	return err
}

func GetUser(db *sql.DB, login string) (string, int, error) {
	query := "SELECT password, is_admin FROM Users WHERE login = ?"
	var password string
	var is_admin int

	err := db.QueryRow(query, login).Scan(&password, &is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return "", -1, err
	}

	return password, is_admin, nil
}
