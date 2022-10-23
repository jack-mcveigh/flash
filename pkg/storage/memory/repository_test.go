package memory

import (
	"reflect"
	"testing"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
)

func TestAddCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		card    adding.Card
		want    []Card
		wantErr error
	}{
		{
			name: "Normal",
			card: adding.Card{Title: "Subject1", Desc: "Value1"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
			},
			wantErr: nil,
		},
		{
			name: "Empty Title",
			card: adding.Card{Title: "", Desc: "Value1"},
			want: []Card{
				{Title: "", Desc: "Value1"},
			},
			wantErr: nil,
		},
		{
			name: "Empty Desc",
			card: adding.Card{Title: "Subject1", Desc: ""},
			want: []Card{
				{Title: "Subject1", Desc: ""},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			err := r.AddCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
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
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			for _, c := range tt.want {
				r.AddCard(adding.Card{Title: c.Title, Desc: c.Desc})
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestDeleteCardSingle(t *testing.T) {
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
		{
			name: "Card not found",
			card: deleting.Card{Title: "Subject4"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
			wantErr: ErrCardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			}
			err := r.DeleteCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestDeleteCardMultiple(t *testing.T) {
	tests := []struct {
		name  string
		cards []deleting.Card
		want  []Card
	}{
		{
			name: "Delete All",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			want: []Card{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			for _, c := range tt.cards {
				r.cards = append(r.cards, Card{Title: c.Title, Desc: ""})
			}

			for _, c := range tt.cards {
				r.DeleteCard(c)
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
			r := New()
			for _, c := range tt.want {
				r.cards = append(r.cards, Card{Title: c.Title, Desc: c.Desc})
			}
			cards := r.GetCards()

			if !reflect.DeepEqual(tt.want, cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, cards)
			}
		})
	}
}
