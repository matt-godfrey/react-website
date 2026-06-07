package users

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username string, email string, passwordHash string, isActive bool, createdAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (username, email, password_hash, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, username, email, passwordHash, isActive, createdAt)

	if isUniqueViolation(err) {
		return ErrUserAlreadyExists
	}
	return err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, is_active, created_at
		FROM users
		WHERE email = $1
	`, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, is_active, created_at
		FROM users
		WHERE id = $1
	`, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsActive, &user.CreatedAt)
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
