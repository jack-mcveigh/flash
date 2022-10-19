package deleting

type Service interface {
	DeleteCard(Card)
}

type Repository interface {
	DeleteCard(Card)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) DeleteCard(c Card) {
	s.r.DeleteCard(c)
}
