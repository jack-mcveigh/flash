package memory

import (
	"reflect"
	"testing"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
)

func newRepositoryWithCards() *Repository {
	r := NewRepository()
	r.cards = append(
		r.cards,
		Card{Title: "Subject1", Desc: "Value1"},
		Card{Title: "Subject2", Desc: "Value2"},
		Card{Title: "Subject3", Desc: "Value3"},
	)
	return r
}

func TestAddCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		card    adding.Card
		want    []Card
		wantErr error
	}{
		{
			name: "Normal",
			card: adding.Card{Title: "Subject4", Desc: "Value4"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
				{Title: "Subject4", Desc: "Value4"},
			},
			wantErr: nil,
		},
		{
			name: "Empty Title",
			card: adding.Card{Title: "", Desc: "Value4"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
				{Title: "", Desc: "Value4"},
			},
			wantErr: nil,
		},
		{
			name: "Empty Desc",
			card: adding.Card{Title: "Subject4", Desc: ""},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
				{Title: "Subject4", Desc: ""},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRepositoryWithCards()
			if err := r.AddCard(tt.card); err != nil {
				if tt.wantErr != err {
					t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
				}
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestAddCardMultiple(t *testing.T) {
	tests := []struct {
		name    string
		cards   []adding.Card
		want    []Card
		wantErr error
	}{
		{
			name: "Normal",
			cards: []adding.Card{
				{Title: "Subject4", Desc: "Value4"},
				{Title: "Subject5", Desc: "Value5"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
				{Title: "Subject4", Desc: "Value4"},
				{Title: "Subject5", Desc: "Value5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRepositoryWithCards()
			for _, c := range tt.cards {
				r.AddCard(c)
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		name    string
		card    deleting.Card
		want    []Card
		wantErr error
	}{
		{
			name: "Delete First",
			card: deleting.Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
			wantErr: nil,
		},
		{
			name: "Delete Last",
			card: deleting.Card{Title: "Subject3"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRepositoryWithCards()
			if err := r.DeleteCard(tt.card); err != nil {
				if tt.wantErr != err {
					t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
				}
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestGetCards(t *testing.T) {
	tests := []struct {
		name string
		want []getting.Card
	}{
		{
			name: "Single Card",
			want: []getting.Card{
				{Title: "Subject1", Desc: "Value1"},
			},
		},
		{
			name: "Multiple Cards",
			want: []getting.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository()
			for _, c := range tt.want {
				r.cards = append(r.cards, Card{Title: c.Title, Desc: c.Desc})
			}
			cards := r.GetCards()

			if !reflect.DeepEqual(tt.want, cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}
