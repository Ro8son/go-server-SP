-- name: CreateUser :exec
	INSERT INTO Users (
    login, password, email
  ) VALUES(
    ?, ?, ?
  );

-- name: GetUser :one
  SELECT id, password, is_admin FROM Users 
  WHERE login = ? LIMIT 1;

-- name: AddFile :one
  INSERT INTO Files (
    owner_id, file_name, title, description, coordinates
  ) VALUES(
    ?, ?, ?, ?, ?
  ) RETURNING id;

-- name: GetFiles :many
  SELECT id, file_name FROM Files 
  WHERE owner_id = ?;
