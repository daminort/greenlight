package models

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Movies *MovieService
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Movies: &MovieService{DB: db},
	}
}
