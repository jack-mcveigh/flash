package deleting

import (
	"errors"
	"reflect"
	"testing"
)

var ErrCardNotFound error = errors.New("Card not found")

type repositoryStub struct {
	cards []Card
}

func newRepositoryStub() *repositoryStub {
	r := &repositoryStub{}
	r.cards = []Card{
		{Title: "Subject1"},
		{Title: "Subject2"},
		{Title: "Subject3"},
	}
	return r
}

func (r *repositoryStub) DeleteCard(c Card) {
	index := -1
	for i, card := range r.cards {
		if c.Title == card.Title {
			index = i
		}
	}

	if index < 0 {
		return
	}

	r.cards = append(r.cards[:index], r.cards[index+1:]...)
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		name string
		card Card
		want []Card
	}{
		{
			name: "Normal",
			card: Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepositoryStub()
			ds := NewService(repo)
			ds.DeleteCard(tt.card)

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
