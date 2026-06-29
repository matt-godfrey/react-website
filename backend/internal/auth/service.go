package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matt-godfrey/react-website/internal/mailer"
	"github.com/matt-godfrey/react-website/internal/queue"
	"github.com/matt-godfrey/react-website/internal/sessions"
	"github.com/matt-godfrey/react-website/internal/users"
	"github.com/matt-godfrey/react-website/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string, email string, passwordHash string, isActive bool, createdAt time.Time) (*users.User, error)
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
	LogoutUser(ctx context.Context, sessionId string) error
	Authenticate(ctx context.Context, sessionId string) (*users.User, error)
	CreateSession(ctx context.Context, userId int64, sessionId string) error
}

type svc struct {
	usersRepo    UserRepository
	sessionsRepo SessionRepository
	mailer       mailer.Mailer
	rabbitmq     *queue.RabbitClient
}

func NewService(usersRepo UserRepository, sessionsRepo SessionRepository, mailer mailer.Mailer, rabbitmq *queue.RabbitClient) Service {
	return &svc{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
		mailer:       mailer,
		rabbitmq:     rabbitmq,
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

	user, err := s.usersRepo.CreateUser(ctx, username, email, hashedPassword, isActive, createdAt)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = s.rabbitmq.Send(ctx, "email", "email", amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte("Register Successful"),
		Headers:      amqp.Table{"user_id": user.ID, "username": username, "email": email},
	})
	if err != nil {
		return err
	}

	// // Send mail
	// msg := mailer.Message{
	// 	From:    os.Getenv("RESEND_FROM"),
	// 	To:      []string{user.Email},
	// 	Subject: "Login Successful",
	// 	HTML:    "<p>You have successfully logged in to your account.</p>",
	// 	Text:    "You have successfully logged in to your account.",
	// }
	// _, err = s.mailer.SendMail(ctx, msg)
	// if err != nil {
	// 	return err
	// }

	return nil
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

	return sessionId, nil
}

// LogoutUser handles the logout request
func (s *svc) LogoutUser(ctx context.Context, sessionId string) error {

	err := s.sessionsRepo.Delete(ctx, sessionId)
	if err != nil {
		return err
	}
	return nil
}

// Authenticate validates the session id and returns the user if valid
func (s *svc) Authenticate(ctx context.Context, sessionId string) (*users.User, error) {
	session, err := s.sessionsRepo.FindByID(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	// check if session has expired
	// if so, delete it and return an error
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
