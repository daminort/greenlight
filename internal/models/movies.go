package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
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

type MovieService struct {
	DB *sql.DB
}

var ErrInvalidRuntime = errors.New("invalid runtime format (expected 'N mins')")

// Movie

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

// MovieService

func (ms *MovieService) GetMovies() ([]Movie, error) {
	return nil, nil
}

func (ms *MovieService) GetMovie(id int64) (*Movie, error) {
	return nil, nil
}

func (ms *MovieService) InsertMovie(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return ms.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (ms *MovieService) UpdateMovie(movie *Movie) (*Movie, error) {
	return nil, nil
}

func (ms *MovieService) DeleteMovie(id int64) error {
	return nil
}
