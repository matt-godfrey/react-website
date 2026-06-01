package auth

import (
	"context"

	"github.com/matt-godfrey/react-website/internal/users"
	"github.com/matt-godfrey/react-website/internal/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string, email string, passwordHash string) error
	FindUserByEmail(ctx context.Context, email string) (*users.User, error)
}

type Service interface {
	RegisterUser(ctx context.Context, username string, email string, passwordHash string) error
	LoginUser(ctx context.Context, email string, passwordHash string) (*users.User, error)
}

type svc struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) RegisterUser(ctx context.Context, username string, email string, password string) error {

	// TODO
	// password and email verification

	// hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, username, email, hashedPassword)
	return err
}

func (s *svc) LoginUser(ctx context.Context, email string, password string) (*users.User, error) {

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	// check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, utils.ErrInvalidPassword
	}

	return user, nil
}
