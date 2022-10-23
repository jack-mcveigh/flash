package memory

import (
	"errors"
	"time"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
)

var (
	ErrCardAlreadyExists = errors.New("Card already exists")
	ErrCardNotFound      = errors.New("Card not found")
)

type Clock interface {
	Now() time.Time
}

type clock struct {
}

func (c *clock) Now() time.Time {
	return time.Now()
}

type repository struct {
	cards []Card
	clock Clock
}

func New() *repository {
	r := &repository{}
	r.clock = &clock{}
	return r
}

func (r *repository) AddCard(c adding.Card) error {
	for _, card := range r.cards {
		if card.Title == c.Title {
			return ErrCardAlreadyExists
		}
	}

	t := r.clock.Now()
	r.cards = append(
		r.cards,
		Card{
			Title:   c.Title,
			Desc:    c.Desc,
			Created: t,
			Updated: t,
		},
	)
	return nil
}

func (r *repository) DeleteCard(c deleting.Card) error {
	index := -1
	for i, card := range r.cards {
		if c.Title == card.Title {
			index = i
		}
	}

	if index == -1 {
		return ErrCardNotFound
	}

	r.cards = append(r.cards[:index], r.cards[index+1:]...)
	return nil
}

func (r *repository) GetCards() ([]getting.Card, error) {
	var cards []getting.Card
	for _, c := range r.cards {
		cards = append(cards, getting.Card{
			Title: c.Title,
			Desc:  c.Desc,
		})
	}
	return cards, nil
}

func (r *repository) UpdateCard(c updating.Card) error {
	for i := range r.cards {
		if r.cards[i].Title == c.Title {
			r.cards[i].Desc = c.Desc
			r.cards[i].Updated = r.clock.Now()
			return nil
		}
	}
	return ErrCardNotFound
}
