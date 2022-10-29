package json

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	"github.com/jmcveigh55/flash/pkg/storage"
	scribble "github.com/nanobox-io/golang-scribble"
)

const cardCollection = "card"

var (
	dataPath = "/tmp/.flash"

	ErrCardFound     = errors.New("card already exists")
	ErrCardNotFound  = errors.New("card not found")
	ErrGroupNotFound = errors.New("group not found")
)

func joinSubCollectionPath(c, s string) string {
	if s != "" {
		return c + "/" + strings.Replace(s, ".", "/", -1)
	}
	return c
}

type dbDriver interface {
	Write(string, string, any) error
	Read(string, string, any) error
	ReadAll(string) ([]string, error)
	Delete(string, string) error
}

type repository struct {
	db    dbDriver
	clock storage.Clock
}

func New() (*repository, error) {
	usr, err := user.Current()
	if err == nil {
		dataPath = usr.HomeDir + "/.flash"
	}
	db, err := scribble.New(dataPath, nil)
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
	subCollection := joinSubCollectionPath(cardCollection, g)
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
	subCollection := joinSubCollectionPath(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return ErrGroupNotFound
	}
	if ok := r.checkCardExists(subCollection, c.Title); !ok {
		return ErrCardNotFound
	}

	return r.db.Delete(subCollection, c.Title)
}

func getCardsFromGroup(g string) ([]Card, error) {
	p := dataPath + "/" + g

	cards := []Card{}
	f, err := os.Open(p)
	if err != nil {
		return cards, ErrGroupNotFound
	}
	defer f.Close()

	items, err := f.ReadDir(0)
	if err != nil {
		return cards, err
	}

	for _, i := range items {
		if i.IsDir() || filepath.Ext(i.Name()) != ".json" {
			continue
		}

		jsonFile, err := os.Open(filepath.Join(p, i.Name()))
		if err != nil {
			return cards, err
		}
		defer jsonFile.Close()

		var c Card
		err = json.NewDecoder(jsonFile).Decode(&c)
		if err != nil {
			return cards, err
		}

		cards = append(cards, c)
	}
	return cards, nil
}

func (r *repository) GetCards(g string) ([]getting.Card, error) {
	cards := []getting.Card{}
	subCollection := joinSubCollectionPath(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return cards, ErrGroupNotFound
	}

	cs, err := getCardsFromGroup(subCollection)
	if err != nil {
		return cards, err
	}

	for _, c := range cs {
		cards = append(cards, getting.Card{Title: c.Title, Desc: c.Desc})
	}
	return cards, nil
}

func (r *repository) UpdateCard(g string, c updating.Card) error {
	subCollection := joinSubCollectionPath(cardCollection, g)
	if ok := r.checkGroupExists(subCollection); !ok {
		return ErrGroupNotFound
	}

	if ok := r.checkCardExists(subCollection, c.Title); !ok {
		return ErrCardNotFound
	}

	card := &Card{}
	if err := r.db.Read(subCollection, c.Title, &card); err != nil {
		return err
	}

	u := Card{
		Title:   card.Title,
		Desc:    c.Desc,
		Created: card.Created,
		Updated: r.clock.Now(),
	}
	return r.db.Write(subCollection, card.Title, u)
}
