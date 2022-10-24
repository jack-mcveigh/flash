package updating

import "errors"

var ErrCardEmptyTitle error = errors.New("Card has an empty title")
var ErrCardNotFound error = errors.New("Card not found")

type Service interface {
	UpdateCard(*Card) error
}

type Repository interface {
	UpdateCard(*Card) error
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) UpdateCard(c *Card) error {
	if c.Title == "" {
		return ErrCardEmptyTitle
	}
	return s.r.UpdateCard(c)
}
