package auth

import (
	"context"

	"github.com/matt-godfrey/react-website/internal/users"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string, email string, passwordHash string) error
	FindUserByEmail(ctx context.Context, email string) (*users.User, error)
}

type Service interface {
	RegisterUser(ctx context.Context, username string, email string, passwordHash string) error
}

type svc struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) RegisterUser(ctx context.Context, username string, email string, passwordHash string) error {

	err := s.repo.CreateUser(ctx, username, email, passwordHash)
	return err
}
