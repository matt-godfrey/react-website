package sessions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertSession(ctx context.Context, session Session) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)`,
		session.Id, session.UserId, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *Repository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionID)
	return err
}

func (r *Repository) FindByID(ctx context.Context, sessionID string) (*Session, error) {
	var session Session
	err := r.db.QueryRow(ctx, `SELECT * FROM sessions WHERE id = $1`, sessionID).Scan(
		&session.Id, &session.UserId, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
