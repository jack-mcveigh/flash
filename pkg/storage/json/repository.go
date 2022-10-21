package json

import (
	"encoding/json"
	"errors"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	scribble "github.com/nanobox-io/golang-scribble"
)

const (
	dataPath       = "/tmp/flash"
	CardCollection = "card"
)

var ErrCardNotFound error = errors.New("Card not found")

type Repository struct {
	db *scribble.Driver
}

func New() (*Repository, error) {
	db, err := scribble.New(dataPath, nil)
	return &Repository{db}, err
}

func (r *Repository) AddCard(c adding.Card) error {
	return nil
}

func (r *Repository) DeleteCard(c deleting.Card) error {
	return nil
}

func (r *Repository) GetCards() []getting.Card {
	cards := []getting.Card{}
	records, err := r.db.ReadAll(CardCollection)
	if err != nil {
		// TODO: Add error handling
		return cards
	}

	for _, r := range records {
		var c Card
		if err := json.Unmarshal([]byte(r), &c); err != nil {
			// TODO: Add error handling
			return cards
		}

		cards = append(cards, getting.Card{Title: c.Title, Desc: c.Desc})
	}
	return cards
}
