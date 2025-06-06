// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const addFile = `-- name: AddFile :one
  INSERT INTO Files (
    owner_id, file_name, title, description, coordinates
  ) VALUES(
    ?, ?, ?, ?, ?
  ) RETURNING id
`

type AddFileParams struct {
	OwnerID     int64          `json:"owner_id"`
	FileName    string         `json:"file_name"`
	Title       sql.NullString `json:"title"`
	Description sql.NullString `json:"description"`
	Coordinates sql.NullString `json:"coordinates"`
}

func (q *Queries) AddFile(ctx context.Context, arg AddFileParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, addFile,
		arg.OwnerID,
		arg.FileName,
		arg.Title,
		arg.Description,
		arg.Coordinates,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :exec
	INSERT INTO Users (
    login, password, email
  ) VALUES(
    ?, ?, ?
  )
`

type CreateUserParams struct {
	Login    string         `json:"login"`
	Password string         `json:"password"`
	Email    sql.NullString `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.Login, arg.Password, arg.Email)
	return err
}

const getFiles = `-- name: GetFiles :many
  SELECT id, file_name FROM Files 
  WHERE owner_id = ?
`

type GetFilesRow struct {
	ID       int64  `json:"id"`
	FileName string `json:"file_name"`
}

func (q *Queries) GetFiles(ctx context.Context, ownerID int64) ([]GetFilesRow, error) {
	rows, err := q.db.QueryContext(ctx, getFiles, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFilesRow
	for rows.Next() {
		var i GetFilesRow
		if err := rows.Scan(&i.ID, &i.FileName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
  SELECT id, password, is_admin FROM Users 
  WHERE login = ? LIMIT 1
`

type GetUserRow struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
	IsAdmin  int64  `json:"is_admin"`
}

func (q *Queries) GetUser(ctx context.Context, login string) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, login)
	var i GetUserRow
	err := row.Scan(&i.ID, &i.Password, &i.IsAdmin)
	return i, err
}
