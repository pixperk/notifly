package user

import (
	"database/sql"

	database "github.com/pixperk/notifly/user/db/sqlc"
)

type Store struct {
	*database.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: database.New(db),
		db:      db,
	}
}
