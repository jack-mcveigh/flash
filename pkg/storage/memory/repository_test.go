package memory

import (
	"reflect"
	"testing"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
)

func TestAddCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		card    adding.Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			card:    adding.Card{Title: "Subject1", Desc: "Value1"},
			want:    []Card{{Title: "Subject1", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name:    "Empty Title",
			card:    adding.Card{Title: "", Desc: "Value1"},
			want:    []Card{{Title: "", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name:    "Empty Desc",
			card:    adding.Card{Title: "Subject1", Desc: ""},
			want:    []Card{{Title: "Subject1", Desc: ""}},
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
			cards: []adding.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name: "Duplicate Title",
			cards: []adding.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject1", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
			},
			wantErr: ErrCardAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			var err error
			for _, c := range tt.cards {
				err = r.AddCard(c)
			}

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
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
			name: "Normal",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			want: []Card{},
		},
		{
			name: "One Card Not Found",
			cards: []deleting.Card{
				{Title: "Subject2"},
				{Title: "Subject3"},
				{Title: "Subject4"},
			},
			want: []Card{{Title: "Subject1", Desc: "Value1"}},
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
		name    string
		want    []getting.Card
		wantErr error
	}{
		{
			name:    "Single Card",
			want:    []getting.Card{{Title: "Subject1", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name: "Multiple Cards",
			want: []getting.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			for _, c := range tt.want {
				r.cards = append(r.cards, Card{Title: c.Title, Desc: c.Desc})
			}
			cards, err := r.GetCards()
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, cards)
			}
		})
	}
}

func TestUpdateCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		card    updating.Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			card:    updating.Card{Title: "Subject1", Desc: "Value2"},
			want:    []Card{{Title: "Subject1", Desc: "Value2"}},
			wantErr: nil,
		},
		{
			name:    "Empty Title",
			card:    updating.Card{Title: "Subject2", Desc: "Value2"},
			want:    []Card{{Title: "Subject1", Desc: "Value1"}},
			wantErr: ErrCardNotFound,
		},
		{
			name:    "Empty Desc",
			card:    updating.Card{Title: "Subject1", Desc: ""},
			want:    []Card{{Title: "Subject1", Desc: ""}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.cards = []Card{{Title: "Subject1", Desc: "Value1"}}
			err := r.UpdateCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}

func TestUpdateCardMultiple(t *testing.T) {
	tests := []struct {
		name  string
		cards []updating.Card
		want  []Card
	}{
		{
			name: "Normal",
			cards: []updating.Card{
				{Title: "Subject1", Desc: "Value3"},
				{Title: "Subject2", Desc: "Value4"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value3"},
				{Title: "Subject2", Desc: "Value4"},
			},
		},
		{
			name: "One Card Not Found",
			cards: []updating.Card{
				{Title: "Subject2", Desc: "Value3"},
				{Title: "Subject3", Desc: "Value4"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			r.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			}

			for _, c := range tt.cards {
				r.UpdateCard(c)
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}
