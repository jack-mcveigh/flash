package getting

import (
	"reflect"
	"testing"
)

type repositoryStub struct {
	cards []Card
}

func (r *repositoryStub) GetCards() []Card {
	return r.cards
}

func TestGetCards(t *testing.T) {
	tests := []struct {
		name string
		want []Card
	}{
		{
			name: "Normal",
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
		},
		{
			name: "Single Card",
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
			},
		},
		{
			name: "No Cards",
			want: []Card{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			repo.cards = tt.want
			gs := NewService(repo)
			got := gs.GetCards()

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, got)
			}
		})
	}
}
