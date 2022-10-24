package deleting

import (
	"reflect"
	"testing"
)

type repositoryStub struct {
	cards []Card
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
		return ErrCardNotFound
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
				{Title: "Group.Subject2"},
				{Title: "Group.Subject3"},
			},
			wantErr: nil,
		},
		{
			name:  "No Group",
			group: "",
			card:  Card{Title: "Subject1"},
			want: []Card{
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.Subject3"},
			},
			wantErr: nil,
		},
		{
			name:  "Card not found",
			group: "Group",
			card:  Card{Title: "Subject4"},
			want: []Card{
				{Title: "Subject1"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.Subject3"},
			},
			wantErr: ErrCardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			repo.cards = []Card{
				{Title: "Subject1"},
				{Title: "Group.Subject1"},
				{Title: "Group.Subject2"},
				{Title: "Group.Subject3"},
			}
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
