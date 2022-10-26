package memory

import (
	"reflect"
	"testing"
	"time"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
)

type clockStub struct{}

func (c *clockStub) Now() time.Time {
	return time.Time{}
}

func TestAddCardSingle(t *testing.T) {
	tests := []struct {
		name    string
		group   string
		card    adding.Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			group:   "Group",
			card:    adding.Card{Title: "Subject1", Desc: "Value1"},
			want:    []Card{{Title: "Group.Subject1", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name:    "Empty Title",
			group:   "Group",
			card:    adding.Card{Title: "", Desc: "Value1"},
			want:    []Card{{Title: "Group.", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name:    "Empty Desc",
			group:   "Group",
			card:    adding.Card{Title: "Subject1", Desc: ""},
			want:    []Card{{Title: "Group.Subject1", Desc: ""}},
			wantErr: nil,
		},
		{
			name:    "Empty Group",
			group:   "",
			card:    adding.Card{Title: "Subject1", Desc: ""},
			want:    []Card{{Title: "Subject1", Desc: ""}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			err := r.AddCard(tt.group, tt.card)

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
		group   string
		cards   []adding.Card
		want    []Card
		wantErr error
	}{
		{
			name:  "Normal",
			group: "Group",
			cards: []adding.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Duplicate Title",
			group: "Group",
			cards: []adding.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject1", Desc: "Value2"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
			},
			wantErr: ErrCardFound,
		},
		{
			name:  "Empty Group",
			group: "",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			var err error
			for _, c := range tt.cards {
				err = r.AddCard(tt.group, c)
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
		group   string
		card    deleting.Card
		want    []Card
		wantErr error
	}{
		{
			name:  "Delete Second",
			group: "Group",
			card:  deleting.Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			},
			wantErr: nil,
		},
		{
			name:  "Delete Last",
			group: "Group",
			card:  deleting.Card{Title: "Subject3"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Card not found",
			group: "Group",
			card:  deleting.Card{Title: "Subject4"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			},
			wantErr: ErrCardNotFound,
		},
		{
			name:  "Empty Group",
			group: "",
			card:  deleting.Card{Title: "Subject1"},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			r.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			}
			err := r.DeleteCard(tt.group, tt.card)
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
		group string
		cards []deleting.Card
		want  []Card
	}{
		{
			name:  "Normal",
			group: "Group",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			want: []Card{{Title: "Subject1", Desc: "Value1"}},
		},
		{
			name:  "One Card Not Found",
			group: "Group",
			cards: []deleting.Card{
				{Title: "Subject2"},
				{Title: "Subject3"},
				{Title: "Subject4"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
			},
		},
		{
			name:  "Empty Group",
			group: "",
			cards: []deleting.Card{
				{Title: "Subject1"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			r.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			}

			for _, c := range tt.cards {
				r.DeleteCard(tt.group, c)
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
		group   string
		want    []getting.Card
		wantErr error
	}{
		{
			name:    "Single Card",
			group:   "Group",
			want:    []getting.Card{{Title: "Group.Subject1", Desc: "Value1"}},
			wantErr: nil,
		},
		{
			name:  "Multiple Cards",
			group: "Group",
			want: []getting.Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.Subject3", Desc: "Value3"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			for _, c := range tt.want {
				r.cards = append(r.cards, Card{Title: c.Title, Desc: c.Desc})
			}
			cards, err := r.GetCards(tt.group)

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
		group   string
		card    updating.Card
		want    []Card
		wantErr error
	}{
		{
			name:    "Normal",
			group:   "Group",
			card:    updating.Card{Title: "Subject1", Desc: "Value2"},
			want:    []Card{{Title: "Group.Subject1", Desc: "Value2"}},
			wantErr: nil,
		},
		{
			name:    "Empty Desc",
			group:   "Group",
			card:    updating.Card{Title: "Subject1", Desc: ""},
			want:    []Card{{Title: "Group.Subject1", Desc: ""}},
			wantErr: nil,
		},
		{
			name:    "Card Not Found",
			group:   "Group",
			card:    updating.Card{Title: "Subject2", Desc: "Value2"},
			want:    []Card{{Title: "Group.Subject1", Desc: "Value1"}},
			wantErr: ErrCardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			r.cards = []Card{{Title: "Group.Subject1", Desc: "Value1"}}
			err := r.UpdateCard(tt.group, tt.card)

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
		group string
		cards []updating.Card
		want  []Card
	}{
		{
			name:  "Normal",
			group: "Group",
			cards: []updating.Card{
				{Title: "Group.Subject1", Desc: "Value3"},
				{Title: "Group.Subject2", Desc: "Value4"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value3"},
				{Title: "Group.Subject2", Desc: "Value4"},
			},
		},
		{
			name:  "One Card Not Found",
			group: "Group",
			cards: []updating.Card{
				{Title: "Group.Subject2", Desc: "Value3"},
				{Title: "Group.Subject3", Desc: "Value4"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{}
			r.clock = &clockStub{}
			r.cards = []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
			}

			for _, c := range tt.cards {
				r.UpdateCard("", c)
			}

			if !reflect.DeepEqual(tt.want, r.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, r.cards)
			}
		})
	}
}
