package user

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"server/database"
)

func AddUser(db *sql.DB, login, passwordHash, email string) (string, error) {
	// Sanitize user input
	login = strings.Replace(login, "/", "∕", -1)
	login = strings.TrimSpace(login)

	if login == "" || passwordHash == "" {
		return "", errors.New("Empty login data")
	}

	err := database.AddUser(db, login, passwordHash, email)
	if err != nil {
		return "", err
	}

	// Create user directory
	userdir := filepath.Join("../storage/users/", login)
	if _, err := os.Stat(userdir); !os.IsNotExist(err) {
		return login, nil
	}

	err = os.MkdirAll(userdir, 0755)
	if err != nil {
		return login, err
	}

	return login, nil
}
