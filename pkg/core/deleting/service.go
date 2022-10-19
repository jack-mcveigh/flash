package deleting

type Service interface {
	DeleteCard(Card) error
}

type Repository interface {
	DeleteCard(Card) error
}

type service struct {
	r Repository
}

func New(r Repository) *service {
	return &service{r}
}

func (s *service) DeleteCard(c Card) error {
	return s.r.DeleteCard(c)
}
