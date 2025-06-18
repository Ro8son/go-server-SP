-- name: CreateUser :exec
INSERT INTO users (
  login, password, email
) VALUES(
  ?, ?, ?
);

-- name: UpdateUser :one
UPDATE users
SET email = ?, profile = ?
WHERE id = ?
RETURNING *;

-- name: ChangeRole :one
UPDATE users
SET is_admin = ?
WHERE id = ?
RETURNING *;

-- name: GetUser :one
SELECT id, login, password, is_admin FROM users 
WHERE id = ? LIMIT 1;

-- name: GetUserByLogin :one
SELECT id, login, password, is_admin FROM users 
WHERE login = ? LIMIT 1;

-- name: GetLogin :one
SELECT login FROM users 
WHERE id = ? LIMIT 1;

-- name: GetPassword :one
SELECT password FROM users 
WHERE id = ? LIMIT 1;

-- name: GetRole :one
SELECT is_admin FROM users 
WHERE id = ? LIMIT 1;

-- name: GetProfile :one
SELECT profile FROM users 
WHERE id = ? LIMIT 1;

-- name: GetEmail :one
SELECT email FROM users 
WHERE id = ? LIMIT 1;

-- name: AddFile :one
INSERT INTO files (
  owner_id, file_name, title, description, coordinates
) VALUES(
  ?, ?, ?, ?, ?
) RETURNING id;

-- name: GetFiles :many
SELECT id, file_name FROM files 
WHERE owner_id = ?;

-- name: GetFileOwner :one
SELECT owner_id FROM files
WHERE id = ?;

-- name: AddGuestFile :one
INSERT INTO fileGuestShares (
  file_id, url, expires_at, max_uses
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: GetSharedFiles :many
SELECT fileGuestShares.* FROM fileGuestShares
LEFT JOIN files ON files.id = fileGuestShares.file_id
WHERE files.owner_id = ?;

-- name: GetShareDownload :one
SELECT files.* FROM fileGuestShares
LEFT JOIN files ON files.id = fileGuestShares.file_id
WHERE fileGuestShares.url = ? AND fileGuestShares.id = ?;

-- name: AddAlbum :exec
INSERT INTO album (
  title
) VALUES (
  ?
);

-- name: getAlbums :many
SELECT title FROM album
WHERE user_id;

