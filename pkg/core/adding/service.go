package adding

import "errors"

var ErrEmptyTitle error = errors.New("")

type Service interface {
	AddCard(Card) error
}

type Repository interface {
	AddCard(Card) error
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) AddCard(c Card) error {
	if c.Title == "" {
		return ErrEmptyTitle
	}
	return s.r.AddCard(c)
}
