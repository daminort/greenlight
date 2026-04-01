package users

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"greenlight.damian.net/internal/errors_manager"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, activated, created_at, version
		FROM users
		WHERE email = $1`

	var u User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Activated,
		&u.CreatedAt,
		&u.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsManager.ErrRecordNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) Get(id int64) (*User, error) {
	query := `
		SELECT id, name, email, pwd, activated, created_at, version
		FROM users
		WHERE id = $1`

	var u User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Pwd.Hash,
		&u.Activated,
		&u.CreatedAt,
		&u.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsManager.ErrRecordNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) Create(u *User) error {
	query := `
		INSERT INTO users (name, email, pwd, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{u.Name, u.Email, u.Pwd.Hash, u.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.CreatedAt, &u.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return errorsManager.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *Repository) Update(u *User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, pwd = $4, activated = $5, version = version + 1
		WHERE id = $1 AND version = $6
		RETURNING version`

	args := []any{u.ID, u.Name, u.Email, u.Pwd.Hash, u.Activated, u.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&u.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errorsManager.ErrEditConflict
		case strings.Contains(err.Error(), "users_email_key"):
			return errorsManager.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
