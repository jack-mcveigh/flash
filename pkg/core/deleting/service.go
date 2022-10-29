package deleting

import "errors"

var ErrCardEmptyTitle error = errors.New("Card has an empty title")

type Service interface {
	DeleteCard(string, Card) error
}

type Repository interface {
	DeleteCard(string, Card) error
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) DeleteCard(g string, c Card) error {
	if c.Title == "" {
		return ErrCardEmptyTitle
	}
	return s.r.DeleteCard(g, c)
}
