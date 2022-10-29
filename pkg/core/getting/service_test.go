package getting

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var errGroupNotFound error = errors.New("Group not found")

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

func (r *repositoryStub) GetCards(g string) ([]Card, error) {
	cards := []Card{}
	for _, c := range r.cards {
		if strings.HasPrefix(c.Title, g) {
			cards = append(cards, c)
		}
	}
	if len(cards) == 0 {
		return cards, ErrGroupNotFound
	}
	return cards, nil
}

func TestGetCards(t *testing.T) {
	tests := []struct {
		group   string
		name    string
		want    []Card
		wantErr error
	}{
		{
			name:  "Normal",
			group: "Group",
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "SubGroup",
			group: "Group.SubGroup",
			want: []Card{
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "No Group",
			group: "",
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:    "Group Not Found",
			group:   "NotFound",
			want:    []Card{},
			wantErr: errGroupNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newRepositoryStubWithCards()
			gs := New(repo)
			got, err := gs.GetCards(tt.group)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, got)
			}
		})
	}
}
