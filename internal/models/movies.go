package models

import (
	"fmt"
	"strconv"
	"time"
)

type Runtime int32

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"-"`
}

func (v Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", v)
	quotedValues := strconv.Quote(jsonValue)

	return []byte(quotedValues), nil
}
