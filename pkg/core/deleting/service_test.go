package deleting

import (
	"errors"
	"reflect"
	"testing"
)

var errCardNotFound error = errors.New("card not found")

type repositoryStub struct {
	cards []Card
}

func newRepositoryStubWithCards() *repositoryStub {
	return &repositoryStub{
		cards: []Card{
			{Title: "Subject1"},
			{Title: "Subject2"},
			{Title: "Group.Subject1"},
			{Title: "Group.Subject2"},
			{Title: "Group.SubGroup.Subject1"},
			{Title: "Group.SubGroup.Subject2"},
		},
	}
}

func (r *repositoryStub) DeleteCard(g string, c Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

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
		group   string
		card    Card
		want    []Card
		wantErr error
	}{
		{
			name:  "Normal",
			group: "Group",
			card:  Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Group.Subject2"},
				{Title: "Group.SubGroup.Subject1"},
				{Title: "Group.SubGroup.Subject2"},
			},
			wantErr: nil,
		},
		{
			name:  "Sub Group",
			group: "Group.SubGroup",
			card:  Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.SubGroup.Subject2"},
			},
			wantErr: nil,
		},
		{
			name:  "No Group",
			group: "",
			card:  Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject2"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.SubGroup.Subject1"},
				{Title: "Group.SubGroup.Subject2"},
			},
			wantErr: nil,
		},
		{
			name:  "Card Not Found",
			group: "Group",
			card:  Card{Title: "Subject3"},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.SubGroup.Subject1"},
				{Title: "Group.SubGroup.Subject2"},
			},
			wantErr: errCardNotFound,
		},
		{
			name:  "Card Empty Title",
			group: "Group",
			card:  Card{Title: ""},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.SubGroup.Subject1"},
				{Title: "Group.SubGroup.Subject2"},
			},
			wantErr: ErrCardEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepositoryStubWithCards()
			ds := New(repo)
			err := ds.DeleteCard(tt.group, tt.card)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
