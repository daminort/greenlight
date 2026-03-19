package movies

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"greenlight.damian.net/internal/errors_manager"
	"greenlight.damian.net/internal/pkg/filters"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetList(params GetMoviesParams) ([]Movie, *filters.Meta, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, title, year, runtime, genres, created_at, version
		FROM movies
		WHERE 
		    (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND (genres && $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, params.Filters.SortColumn(), params.Filters.SortDirection())

	args := []any{
		params.Filters.Search,
		pq.Array(params.Genres),
		params.Filters.Limit(),
		params.Filters.Offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, &filters.Meta{}, err
	}
	defer rows.Close()

	movies := []Movie{}
	totalRows := 0
	for rows.Next() {
		var m Movie
		err := rows.Scan(&totalRows, &m.ID, &m.Title, &m.Year, &m.Runtime, pq.Array(&m.Genres), &m.CreatedAt, &m.Version)
		if err != nil {
			return nil, &filters.Meta{}, err
		}

		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, &filters.Meta{}, err
	}

	meta := filters.NewMeta(totalRows, params.Filters.Page, params.Filters.PageSize)

	return movies, meta, nil
}

func (r *Repository) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, errorsManager.ErrRecordNotFound
	}

	query := `
		SELECT id, title, year, runtime, genres, created_at, version
		FROM movies
		WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, id)
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
			return nil, errorsManager.ErrRecordNotFound
		}
		return nil, err
	}

	return &movie, nil
}

func (r *Repository) Create(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (r *Repository) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title = $2, year = $3, runtime = $4, genres = $5, version = version + 1
		WHERE id = $1 AND version = $6
		RETURNING version`

	args := []any{movie.ID, movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorsManager.ErrEditConflict
		}
		return err
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	if id < 1 {
		return errorsManager.ErrRecordNotFound
	}

	query := `
		DELETE FROM movies 
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errorsManager.ErrRecordNotFound
	}

	return nil
}
