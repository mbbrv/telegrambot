package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

var Db *sqlx.DB

type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type TimeAt struct {
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
