// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (identifier, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, identifier, name, password_hash, created_at
`

type CreateUserParams struct {
	Identifier   string `json:"identifier"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*Users, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Identifier, arg.Name, arg.PasswordHash)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Identifier,
		&i.Name,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return &i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, identifier, name, password_hash, created_at FROM users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (*Users, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Identifier,
		&i.Name,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return &i, err
}

const getUserByIdentifier = `-- name: GetUserByIdentifier :one
SELECT id, identifier, name, password_hash, created_at FROM users
WHERE identifier = $1
`

func (q *Queries) GetUserByIdentifier(ctx context.Context, identifier string) (*Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByIdentifier, identifier)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Identifier,
		&i.Name,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return &i, err
}
