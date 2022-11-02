package json

import (
	"encoding/json"
	"errors"
	"os/user"
	"strings"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	"github.com/jmcveigh55/flash/pkg/storage"
)

const cardCollection = "card"

var (
	dataPath = "/tmp/.flash"

	ErrCardFound     = errors.New("card already exists")
	ErrCardNotFound  = errors.New("card not found")
	ErrGroupNotFound = errors.New("group not found")
)

func joinCollectionPaths(c1, c2 string) string {
	if c1 == "" {
		return c2
	}
	return c1 + "/" + strings.Replace(c2, ".", "/", -1)
}

type repository struct {
	db    DBDriver
	clock storage.Clock
}

func New() (*repository, error) {
	usr, err := user.Current()
	if err == nil {
		dataPath = usr.HomeDir + "/.flash"
	}
	db, err := newDb(dataPath)
	c := storage.NewClock()
	return &repository{db, c}, err
}

func (r *repository) checkCardExists(coll, title string) bool {
	if err := r.db.Read(coll, title, &Card{}); err != nil {
		return false
	}
	return true
}

func (r *repository) checkGroupExists(g string) bool {
	if _, err := r.db.ReadAll(g); err != nil {
		return false
	}
	return true
}

func (r *repository) AddCard(g string, c adding.Card) error {
	subCollection := joinCollectionPaths(cardCollection, g)
	if ok := r.checkCardExists(subCollection, c.Title); ok {
		return ErrCardFound
	}

	t := r.clock.Now()
	card := Card{
		Title:   c.Title,
		Desc:    c.Desc,
		Created: t,
		Updated: t,
	}

	err := r.db.Write(subCollection, card.Title, card)
	return err
}

func (r *repository) DeleteCard(g string, c deleting.Card) error {
	subCollection := joinCollectionPaths(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return ErrGroupNotFound
	}
	if ok := r.checkCardExists(subCollection, c.Title); !ok {
		return ErrCardNotFound
	}

	return r.db.Delete(subCollection, c.Title)
}

func (r *repository) GetCards(g string) ([]getting.Card, error) {
	cards := []getting.Card{}
	subCollection := joinCollectionPaths(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return cards, ErrGroupNotFound
	}

	items, err := r.db.ReadAll(subCollection)
	if err != nil {
		return cards, err
	}

	var c getting.Card
	for _, item := range items {
		if err := json.Unmarshal([]byte(item), &c); err != nil {
			return cards, err
		}
		cards = append(cards, getting.Card{Title: c.Title, Desc: c.Desc})
	}
	return cards, nil
}

func (r *repository) UpdateCard(g string, c updating.Card) error {

	subCollection := joinCollectionPaths(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return ErrGroupNotFound
	}

	if ok := r.checkCardExists(subCollection, c.Title); !ok {
		return ErrCardNotFound
	}

	card := Card{}
	if err := r.db.Read(subCollection, c.Title, &card); err != nil {
		return err
	}

	u := Card{
		Title:   c.Title,
		Desc:    c.Desc,
		Created: card.Created,
		Updated: r.clock.Now(),
	}

	return r.db.Write(subCollection, u.Title, u)
}
