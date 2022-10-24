package deleting

import "errors"

var ErrCardNotFound error = errors.New("Card not found")

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
	return s.r.DeleteCard(g, c)
}
