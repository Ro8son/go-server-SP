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

type File struct {
	Id       int64  `json:"id"`
	FileName string `json:"file_name"`
	File     string `json:"file"`
	// checksum
}

func GetFileTitles(db *sql.DB, userId int) ([]File, error) {
	query := `SELECT id, file_name FROM Files WHERE owner_id = ?`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var id int64
		var fileName string
		if err := rows.Scan(&id, &fileName); err != nil {
			//
		}
		files = append(files, File{Id: id, FileName: fileName})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
