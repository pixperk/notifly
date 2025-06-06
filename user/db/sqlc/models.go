// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID    `json:"id"`
	Identifier   string       `json:"identifier"`
	Name         string       `json:"name"`
	PasswordHash string       `json:"password_hash"`
	CreatedAt    sql.NullTime `json:"created_at"`
}
