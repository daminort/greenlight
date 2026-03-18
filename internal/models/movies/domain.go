package movies

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"greenlight.damian.net/internal/pkg/filters"
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

type CreateMoviePayload struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Runtime Runtime  `json:"runtime"`
	Genres  []string `json:"genres"`
}

type UpdateMoviePayload struct {
	Title   *string  `json:"title"`
	Year    *int     `json:"year"`
	Runtime *Runtime `json:"runtime"`
	Genres  []string `json:"genres"`
}

type GetMoviesParams struct {
	Genres []string
	*filters.Filters
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
