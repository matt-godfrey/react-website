package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username string, email string, passwordHash string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
	`, username, email, passwordHash)

	if isUniqueViolation(err) {
		return ErrUserAlreadyExists
	}
	return err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash
		FROM users
		WHERE email = $1
	`, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func notFound(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001"
	}
	return false
}
