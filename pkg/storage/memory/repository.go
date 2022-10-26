package memory

import (
	"errors"
	"strings"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	"github.com/jmcveigh55/flash/pkg/storage"
)

var (
	ErrCardFound     = errors.New("card already exists")
	ErrCardNotFound  = errors.New("card not found")
	ErrGroupNotFound = errors.New("group not found")
)

type repository struct {
	cards []Card
	clock storage.Clock
}

func New() *repository {
	r := &repository{}
	r.clock = storage.NewClock()
	return r
}

func (r *repository) AddCard(g string, c adding.Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

	for _, card := range r.cards {
		if card.Title == c.Title {
			return ErrCardFound
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

func (r *repository) DeleteCard(g string, c deleting.Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

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

func (r *repository) GetCards(g string) ([]getting.Card, error) {
	var cards []getting.Card
	for _, c := range r.cards {
		if strings.HasPrefix(c.Title, g) {
			cards = append(cards, getting.Card{
				Title: c.Title,
				Desc:  c.Desc,
			})
		}
	}

	if len(cards) == 0 {
		return cards, ErrGroupNotFound
	}

	return cards, nil
}

func (r *repository) UpdateCard(g string, c updating.Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

	for i := range r.cards {
		if r.cards[i].Title == c.Title {
			r.cards[i].Desc = c.Desc
			r.cards[i].Updated = r.clock.Now()
			return nil
		}
	}
	return ErrCardNotFound
}
