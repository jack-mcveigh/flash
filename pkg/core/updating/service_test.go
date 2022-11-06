package updating

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
			{Title: "Subject1", Desc: "Value1"},
			{Title: "Subject2", Desc: "Value2"},
			{Title: "Group.Subject1", Desc: "Value1"},
			{Title: "Group.Subject2", Desc: "Value2"},
			{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
			{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
		},
	}
}

func (r *repositoryStub) UpdateCard(g string, c Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

	for i := range r.cards {
		if r.cards[i].Title == c.Title {
			r.cards[i].Desc = c.Desc
			return nil
		}
	}
	return errCardNotFound
}

func TestUpdateCard(t *testing.T) {
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
			card:  Card{Title: "Subject1", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value2"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Empty Desc",
			group: "Group",
			card:  Card{Title: "Subject1", Desc: ""},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: ""},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Sub Group",
			group: "Group.SubGroup",
			card:  Card{Title: "Subject1", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "No Group",
			group: "",
			card:  Card{Title: "Subject1", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value2"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Card Not Found",
			group: "Group",
			card:  Card{Title: "Subject3", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: errCardNotFound,
		},
		{
			name:  "Empty Title",
			group: "Group",
			card:  Card{Title: "", Desc: "Value"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: ErrCardEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepositoryStubWithCards()
			us := New(repo)
			err := us.UpdateCard(tt.group, tt.card)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
