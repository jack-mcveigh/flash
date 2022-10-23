package adding

import (
	"reflect"
	"testing"
)

type repositoryStub struct {
	cards []Card
}

func (r *repositoryStub) AddCard(c Card) error {
	r.cards = append(r.cards, c)
	return nil
}

func TestAddCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		card    Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			card:    Card{Title: "Subject", Desc: "Value"},
			want:    []Card{{Title: "Subject", Desc: "Value"}},
			wantErr: nil,
		},
		{
			name:    "Empty Desc",
			card:    Card{Title: "Subject", Desc: ""},
			want:    []Card{{Title: "Subject", Desc: ""}},
			wantErr: nil,
		},
		{
			name:    "Empty Title",
			card:    Card{Title: "", Desc: "Value"},
			want:    nil,
			wantErr: ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			as := New(repo)
			err := as.AddCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}

func TestAddCardMultiple(t *testing.T) {
	tests := []struct {
		name  string
		cards []Card
		want  []Card
	}{
		{
			name: "Normal",
			cards: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
		},
		{
			name: "One Empty Title",
			cards: []Card{
				{Title: "", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Subject2", Desc: "Value2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			as := New(repo)

			for _, card := range tt.cards {
				as.AddCard(card)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
