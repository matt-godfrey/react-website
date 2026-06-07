package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matt-godfrey/react-website/internal/sessions"
	"github.com/matt-godfrey/react-website/internal/users"
	"github.com/matt-godfrey/react-website/internal/utils"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string, email string, passwordHash string, isActive bool, createdAt time.Time) error
	FindUserByEmail(ctx context.Context, email string) (*users.User, error)
	FindByID(ctx context.Context, id int64) (*users.User, error)
}

type SessionRepository interface {
	InsertSession(ctx context.Context, session sessions.Session) error
	FindByID(ctx context.Context, id string) (*sessions.Session, error)
	Delete(ctx context.Context, id string) error
}

type Service interface {
	RegisterUser(ctx context.Context, username string, email string, passwordHash string) error
	LoginUser(ctx context.Context, email string, passwordHash string) (string, error)
	Authenticate(ctx context.Context, sessionId string) (*users.User, error)
	CreateSession(ctx context.Context, userId int64, sessionId string) error
}

type svc struct {
	usersRepo    UserRepository
	sessionsRepo SessionRepository
}

func NewService(usersRepo UserRepository, sessionsRepo SessionRepository) Service {
	return &svc{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
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

	isActive := true
	createdAt := time.Now()

	err = s.usersRepo.CreateUser(ctx, username, email, hashedPassword, isActive, createdAt)
	return err
}

func (s *svc) LoginUser(ctx context.Context, email string, password string) (string, error) {

	user, err := s.usersRepo.FindUserByEmail(ctx, email)
	if err != nil {
		// log.Println(err)
		return "", err
	}

	// check password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", utils.ErrInvalidPassword
	}

	// generate session id
	sessionId := utils.GenerateSessionId()
	fmt.Println(sessionId)

	// store session in db
	err = s.CreateSession(ctx, user.ID, sessionId)
	if err != nil {
		return "", err
	}
	//
	return sessionId, nil
}

func (s *svc) Authenticate(ctx context.Context, sessionId string) (*users.User, error) {
	session, err := s.sessionsRepo.FindByID(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionsRepo.Delete(ctx, session.Id)
		return nil, errors.New("session expired")
	}

	user, err := s.usersRepo.FindByID(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *svc) CreateSession(ctx context.Context, userId int64, sessionId string) error {
	session := sessions.Session{
		Id:        sessionId,
		UserId:    userId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
	err := s.sessionsRepo.InsertSession(ctx, session)
	return err
}
