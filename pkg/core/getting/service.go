package getting

type Service interface {
	GetCards() []Card
}

type Repository interface {
	GetCards() []Card
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) GetCards() []Card {
	return s.r.GetCards()
}
