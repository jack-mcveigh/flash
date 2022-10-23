package json

import (
	"encoding/json"
	"errors"
	"os/user"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	scribble "github.com/nanobox-io/golang-scribble"
)

const cardCollection = "card"

var (
	dataPath = "/tmp/.flash"

	ErrCardAlreadyExists = errors.New("Card already exists")
	ErrCardNotFound      = errors.New("Card not found")
)

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
	cards, _ := r.GetCards()
	for _, card := range cards {
		if card.Title == c.Title {
			return ErrCardAlreadyExists
		}
	}

	card := Card{
		Title: c.Title,
		Desc:  c.Desc,
	}

	err := r.db.Write(cardCollection, card.Title, card)
	return err
}

func (r *Repository) DeleteCard(c deleting.Card) error {
	cards, _ := r.GetCards()

	index := -1
	for i, card := range cards {
		if c.Title == card.Title {
			index = i
		}
	}

	if index == -1 {
		return ErrCardNotFound
	}

	return r.db.Delete(cardCollection, c.Title)
}

func (r *Repository) GetCards() ([]getting.Card, error) {
	cards := []getting.Card{}
	records, err := r.db.ReadAll(cardCollection)
	if err != nil {
		return cards, err
	}

	for _, r := range records {
		var c Card
		if err := json.Unmarshal([]byte(r), &c); err != nil {
			return cards, err
		}

		cards = append(cards, getting.Card{Title: c.Title, Desc: c.Desc})
	}
	return cards, nil
}

func (r *Repository) UpdateCard(c updating.Card) error {
	cards, _ := r.GetCards()

	for _, card := range cards {
		if card.Title == c.Title {
			u := Card{Title: card.Title, Desc: c.Desc}
			return r.db.Write(cardCollection, card.Title, u)
		}
	}
	return ErrCardNotFound
}
