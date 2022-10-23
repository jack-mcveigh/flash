package json

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
)

type dbDriverStub struct {
	cards []Card
}

func (d *dbDriverStub) Write(collection string, resource string, v any) error {
	switch val := v.(type) {
	case Card:
		d.cards = append(d.cards, val)
	default:
		return errors.New("A card was not passed to dbDriverStub.Write")
	}
	return nil
}

func (d *dbDriverStub) ReadAll(collection string) ([]string, error) {
	var resources []string
	for _, c := range d.cards {
		b, err := json.Marshal(c)
		if err != nil {
			return resources, err
		}
		resources = append(resources, string(b))
	}
	return resources, nil
}

func (d *dbDriverStub) Delete(collection string, resource string) error {
	for i, c := range d.cards {
		if c.Title == resource {
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			return nil
		}
	}
	return ErrCardNotFound
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
			db := &dbDriverStub{}
			r := &Repository{db}
			err := r.AddCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
			db := &dbDriverStub{}
			r := &Repository{db}
			for _, c := range tt.want {
				r.AddCard(adding.Card{Title: c.Title, Desc: c.Desc})
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
			db := &dbDriverStub{}
			r := &Repository{db}
			db.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value3"},
			}
			err := r.DeleteCard(tt.card)
			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
			db := &dbDriverStub{}
			r := &Repository{db}
			for _, c := range tt.cards {
				db.cards = append(db.cards, Card{Title: c.Title, Desc: ""})
			}

			for _, c := range tt.cards {
				r.DeleteCard(c)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
			db := &dbDriverStub{}
			r := &Repository{db}
			for _, c := range tt.want {
				db.cards = append(db.cards, Card{Title: c.Title, Desc: c.Desc})
			}
			cards := r.GetCards()

			if !reflect.DeepEqual(tt.want, cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, cards)
			}
		})
	}
}
