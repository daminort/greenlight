package models

import (
	"context"
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
	query := `
		SELECT id, title, year, runtime, genres, created_at, version
		FROM movies`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := ms.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.Runtime, pq.Array(&m.Genres), &m.CreatedAt, &m.Version)
		if err != nil {
			return nil, err
		}

		movies = append(movies, m)
	}

	return movies, nil
}

func (ms *MovieService) GetMovie(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, year, runtime, genres, created_at, version
		FROM movies
		WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := ms.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.CreatedAt,
		&movie.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &movie, nil
}

func (ms *MovieService) InsertMovie(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return ms.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (ms *MovieService) UpdateMovie(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $2, year = $3, runtime = $4, genres = $5, version = version + 1
		WHERE id = $1 AND version = $6
		RETURNING version`

	args := []any{movie.ID, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ms.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (ms *MovieService) DeleteMovie(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movies 
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := ms.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrRecordNotFound
	}

	return nil
}
