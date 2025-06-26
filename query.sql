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

-- name: GetIsAdmin :one
SELECT is_admin FROM users 
WHERE id = ? LIMIT 1;

-- name: AddFile :one
INSERT INTO files (
  owner_id, file_name, title, description, coordinates, checksum
) VALUES(
  ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: AddTag :one
INSERT INTO tags (
  name
) VALUES (
  ?
)
RETURNING id;

-- name: GetTags :many
select id, name
FROM tags;

-- name: GetTagByName :one
SELECT id, name
FROM tags
WHERE name = ?;

-- name: GetTagById :one
SELECT id, name
FROM tags
WHERE id = ?;

-- name: TagsConnect :exec
INSERT INTO fileTags (
  file_id, tag_id
) VALUES (
  ?, ?
);

-- name: GetTagsByFile :many
SELECT tag_id 
FROM fileTags
WHERE file_id = ?;

-- name: GetFilesByTag :many
SELECT files.* FROM fileTags
LEFT JOIN files ON files.id = fileTags.file_id
WHERE fileTags.tag_id = ? AND files.owner_id = ?;


-- name: GetFiles :many
SELECT id, file_name, checksum, created_at FROM files 
WHERE owner_id = ?;

-- name: GetFile :one
SELECT * FROM files
WHERE id = ?;

-- name: GetFileOwner :one
SELECT owner_id FROM files
WHERE id = ?;

-- name: DeleteFile :exec
DELETE FROM files
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
WHERE (files.owner_id = ?, ? = 1);

-- name: GetShareDownload :one
SELECT files.* FROM fileGuestShares
LEFT JOIN files ON files.id = fileGuestShares.file_id
WHERE fileGuestShares.url = ? AND fileGuestShares.id = ?;

-- name: GetShareUseCount :one
SELECT max_uses FROM fileGuestShares
WHERE id = ?;

-- name: DecrementShareUses :exec
UPDATE fileguestshares
SET max_uses = max_uses - 1
WHERE id = ? AND max_uses > 0;


-- name: AddAlbum :exec
INSERT INTO album (
  title, owner_id, cover_id
) VALUES (
  ?, ?, ?
);

-- name: GetAlbums :many
SELECT * FROM album
WHERE owner_id = ?;

-- name: GetAlbum :one
SELECT * FROM album
WHERE id = ?;

-- name: AddToAlbum :exec
INSERT OR IGNORE INTO fileAlbum (
  file_id, album_id
) VALUES (
  ?, ?
);

-- name: GetFileFromAlbum :many
SELECT file_id
FROM fileAlbum
WHERE album_id = ?;
