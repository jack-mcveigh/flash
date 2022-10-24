package adding

import "errors"

var ErrCardEmptyTitle error = errors.New("Card has an empty title")

type Service interface {
	AddCard(*Card) error
}

type Repository interface {
	AddCard(*Card) error
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) AddCard(c *Card) error {
	if c.Title == "" {
		return ErrCardEmptyTitle
	}
	return s.r.AddCard(c)
}
