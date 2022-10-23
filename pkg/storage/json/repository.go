package json

import (
	"encoding/json"
	"errors"
	"os/user"
	"time"

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

type Clock interface {
	Now() time.Time
}

type clock struct{}

func (c *clock) Now() time.Time {
	return time.Now()
}

type dbDriver interface {
	Write(string, string, any) error
	ReadAll(string) ([]string, error)
	Delete(string, string) error
}

type repository struct {
	db    dbDriver
	clock Clock
}

func New() (*repository, error) {
	usr, err := user.Current()
	if err == nil {
		dataPath = usr.HomeDir + "/.flash"
	}
	db, err := scribble.New(dataPath, nil)
	return &repository{db, &clock{}}, err
}

func (r *repository) AddCard(c adding.Card) error {
	cards, _ := r.GetCards()
	for _, card := range cards {
		if card.Title == c.Title {
			return ErrCardAlreadyExists
		}
	}

	t := r.clock.Now()
	card := Card{
		Title:   c.Title,
		Desc:    c.Desc,
		Created: t,
		Updated: t,
	}

	err := r.db.Write(cardCollection, card.Title, card)
	return err
}

func (r *repository) DeleteCard(c deleting.Card) error {
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

func (r *repository) getCards() ([]Card, error) {
	cards := []Card{}
	records, err := r.db.ReadAll(cardCollection)
	if err != nil {
		return cards, err
	}

	for _, r := range records {
		var c Card
		if err := json.Unmarshal([]byte(r), &c); err != nil {
			return cards, err
		}

		cards = append(cards, c)
	}
	return cards, nil
}

func (r *repository) GetCards() ([]getting.Card, error) {
	cards := []getting.Card{}
	cs, err := r.getCards()
	if err != nil {
		return cards, err
	}

	for _, c := range cs {
		cards = append(cards, getting.Card{Title: c.Title, Desc: c.Desc})
	}
	return cards, nil
}

func (r *repository) UpdateCard(c updating.Card) error {
	cards, _ := r.getCards()

	for _, card := range cards {
		if card.Title == c.Title {
			u := Card{
				Title:   card.Title,
				Desc:    c.Desc,
				Created: card.Created,
				Updated: r.clock.Now(),
			}
			return r.db.Write(cardCollection, card.Title, u)
		}
	}
	return ErrCardNotFound
}
