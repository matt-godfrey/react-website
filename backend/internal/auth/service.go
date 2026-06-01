package auth

import "context"

type Service interface {
	RegisterUser(ctx context.Context) error
}

type svc struct {
	// repo
}

func NewService() Service {
	return &svc{}
}

func (s *svc) RegisterUser(ctx context.Context) error {
	return nil
}
