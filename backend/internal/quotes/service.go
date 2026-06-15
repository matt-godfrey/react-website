package quotes

import (
	"context"
	"math/rand"
)

type QuoteRepository interface {
	FindAllQuotes(ctx context.Context) ([]*Quote, error)
	FindAllQuotesByAuthor(ctx context.Context, author string) ([]*Quote, error)
}

type Service interface {
	GetRandomQuote(ctx context.Context) (*Quote, error)
	GetAllQuotesByAuthor(ctx context.Context, author string) ([]*Quote, error)
}

type svc struct {
	repo QuoteRepository
}

func NewService(repo QuoteRepository) Service {
	return &svc{repo: repo}
}

func (s *svc) GetRandomQuote(ctx context.Context) (*Quote, error) {
	quotes, err := s.repo.FindAllQuotes(ctx)
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return nil, nil
	}
	index := rand.Intn(len(quotes))
	return quotes[index], nil
}

func (s *svc) GetAllQuotesByAuthor(ctx context.Context, author string) ([]*Quote, error) {
	quotes, err := s.repo.FindAllQuotesByAuthor(ctx, author)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}
