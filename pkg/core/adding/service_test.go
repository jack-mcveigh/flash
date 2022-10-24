package adding

import (
	"reflect"
	"testing"
)

type repositoryStub struct {
	cards []Card
}

func (r *repositoryStub) AddCard(g string, c Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}
	r.cards = append(r.cards, c)
	return nil
}

func TestAddCard(t *testing.T) {
	tests := []struct {
		name    string
		group   string
		card    Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			group:   "Group",
			card:    Card{Title: "Subject", Desc: "Value"},
			want:    []Card{{Title: "Group.Subject", Desc: "Value"}},
			wantErr: nil,
		},
		{
			name:    "Empty Desc",
			group:   "Group",
			card:    Card{Title: "Subject", Desc: ""},
			want:    []Card{{Title: "Group.Subject", Desc: ""}},
			wantErr: nil,
		},
		{
			name:    "Empty Group",
			group:   "",
			card:    Card{Title: "Subject", Desc: ""},
			want:    []Card{{Title: "Subject", Desc: ""}},
			wantErr: nil,
		},
		{
			name:    "Empty Title",
			group:   "Group",
			card:    Card{Title: "", Desc: "Value"},
			want:    nil,
			wantErr: ErrCardEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			as := New(repo)
			err := as.AddCard(tt.group, tt.card)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
