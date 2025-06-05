-- name: CreateUser :one
INSERT INTO users (identifier, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByIdentifier :one
SELECT * FROM users
WHERE identifier = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;