package memory

import (
	"errors"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
)

var ErrCardNotFound error = errors.New("Card not found")

type Repository struct {
	cards []Card
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) AddCard(c adding.Card) error {
	r.cards = append(
		r.cards,
		Card{Title: c.Title, Desc: c.Desc},
	)
	return nil
}

func (r *Repository) DeleteCard(c deleting.Card) error {
	index := -1
	for i, card := range r.cards {
		if c.Title == card.Title {
			index = i
		}
	}

	if index < 0 {
		return ErrCardNotFound
	}

	r.cards = append(r.cards[:index], r.cards[index+1:]...)
	return nil
}

func (r *Repository) GetCards() []getting.Card {
	var cards []getting.Card
	for _, c := range r.cards {
		cards = append(cards, getting.Card{
			Title: c.Title,
			Desc:  c.Desc,
		})
	}
	return cards
}
