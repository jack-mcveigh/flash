package json

import (
	"encoding/json"
	"errors"
	"os/user"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	scribble "github.com/nanobox-io/golang-scribble"
)

const cardCollection = "card"

var dataPath = "/tmp/.flash"
var ErrCardNotFound error = errors.New("Card not found")

type dbDriver interface {
	Write(string, string, any) error
	ReadAll(string) ([]string, error)
	Delete(string, string) error
}

type Repository struct {
	db dbDriver
}

func New() (*Repository, error) {
	usr, err := user.Current()
	if err == nil {
		dataPath = usr.HomeDir + "/.flash"
	}
	db, err := scribble.New(dataPath, nil)
	return &Repository{db}, err
}

func (r *Repository) AddCard(c adding.Card) error {
	card := Card{
		Title: c.Title,
		Desc:  c.Desc,
	}

	if err := r.db.Write(cardCollection, card.Title, card); err != nil {
		// TODO: Add error handling
		return err
	}
	return nil
}

func (r *Repository) DeleteCard(c deleting.Card) error {
	// TODO: Add error handling
	return r.db.Delete(cardCollection, c.Title)
}

func (r *Repository) GetCards() []getting.Card {
	cards := []getting.Card{}
	records, err := r.db.ReadAll(cardCollection)
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
