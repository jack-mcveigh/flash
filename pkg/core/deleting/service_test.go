package deleting

import (
	"errors"
	"reflect"
	"testing"
)

var errCardNotFound error = errors.New("Card not found")

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

func (r *repositoryStub) DeleteCard(c Card) error {
	index := -1
	for i, card := range r.cards {
		if c.Title == card.Title {
			index = i
		}
	}

	if index < 0 {
		return errCardNotFound
	}

	r.cards = append(r.cards[:index], r.cards[index+1:]...)
	return nil
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		name    string
		card    Card
		want    []Card
		wantErr error
	}{
		{
			name: "Normal",
			card: Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			wantErr: nil,
		},
		{
			name: "Card not found",
			card: Card{Title: "Subject4"},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			wantErr: errCardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepositoryStub()
			ds := New(repo)
			if err := ds.DeleteCard(tt.card); err != nil {
				if tt.wantErr != err {
					t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
				}
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
