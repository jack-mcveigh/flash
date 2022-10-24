package getting

type Service interface {
	GetCards() ([]*Card, error)
}

type Repository interface {
	GetCards() ([]*Card, error)
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) GetCards() ([]*Card, error) {
	return s.r.GetCards()
}
