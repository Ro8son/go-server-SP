package user

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"server/types"
	"strings"

	"server/database"
)

func AddUser(query *database.Queries, login, passwordHash, email string) (string, error) {
	// Sanitize user input
	login = strings.ReplaceAll(login, "/", "∕")
	login = strings.TrimSpace(login)

	if login == "" || passwordHash == "" {
		return "", errors.New("Empty login data")
	}

	user := database.CreateUserParams{
		Login:    login,
		Password: passwordHash,
		Email:    types.JSONNullString{NullString: sql.NullString{String: email, Valid: email != ""}},
	}

	err := query.CreateUser(context.Background(), user)
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
