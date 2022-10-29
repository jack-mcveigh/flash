package getting

type Service interface {
	GetCards(string) ([]Card, error)
}

type Repository interface {
	GetCards(string) ([]Card, error)
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) GetCards(g string) ([]Card, error) {
	return s.r.GetCards(g)
}
