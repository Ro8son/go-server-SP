package database

import (
	"database/sql"
	"log"
)

func AddUser(db *sql.DB, login, password, email string) error {
	query := "INSERT INTO Users (login, password, email) VALUES(?, ?, ?)"

	_, err := db.Exec(query, login, password, email)
	return err
}

func GetUser(db *sql.DB, login string) (int, string, int, error) {
	query := "SELECT id, password, is_admin FROM Users WHERE login = ?"
	var password string
	var id, is_admin int

	err := db.QueryRow(query, login).Scan(&id, &password, &is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return -1, "", -1, nil
		}
		log.Println("Error retrieving user:", err)
		return -1, "", -1, err
	}

	return id, password, is_admin, nil
}

func AddFile(db *sql.DB, ownerId int, fileName, title, description, coordinates string) (int64, error) {
	query := "INSERT INTO Files (owner_id, file_name, title, description, coordinates) VALUES(?, ?, ?, ?, ?)"

	result, err := db.Exec(query, ownerId, fileName, title, description, coordinates)
	fileId, err := result.LastInsertId()

	return fileId, err
}
