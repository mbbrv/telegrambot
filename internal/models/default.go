package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

type TimeAt struct {
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
