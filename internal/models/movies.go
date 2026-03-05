package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Runtime int

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"-"`
}

var ErrInvalidRuntime = errors.New("invalid runtime format (expected 'N mins')")

func (v Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", v)
	quotedValues := strconv.Quote(jsonValue)

	return []byte(quotedValues), nil
}

func (v *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntime
	}

	parts := strings.Split(unquotedValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntime
	}

	value, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntime
	}

	*v = Runtime(value)

	return nil
}
