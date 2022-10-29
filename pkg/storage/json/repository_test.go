package json

import (
	"encoding/json"
	"errors"
	"io/fs"
	"reflect"
	"strings"
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

type dbDriverStub struct {
	cards []Card
}

func convertCollectionToGroup(coll string) string {
	p := strings.Split(coll, "/")
	return strings.Join(p[1:], ".")
}

func getCardPath(g, t string) string {
	if g == "" {
		return t
	}
	return g + "." + t
}

func (d *dbDriverStub) Write(collection string, resource string, v any) error {
	switch val := v.(type) {
	case Card:
		g := convertCollectionToGroup(collection)
		cardPath := getCardPath(g, val.Title)
		for i := range d.cards {
			if d.cards[i].Title == cardPath {
				d.cards[i].Desc = val.Desc
				return nil
			}
		}
		d.cards = append(d.cards, Card{Title: cardPath, Desc: val.Desc})
		return nil
	default:
		return errors.New("a Card was not passed to dbDriverStub.Write")
	}
}

func (d *dbDriverStub) Read(collection, resource string, v any) error {
	switch val := v.(type) {
	case *Card:
		g := convertCollectionToGroup(collection)
		cardPath := getCardPath(g, resource)
		for _, c := range d.cards {
			if c.Title == cardPath {
				*val = c
				return nil
			}
		}
		return &fs.PathError{}
	default:
		return errors.New("a *Card was not passed to dbDriverStub.Read")
	}
}

func (d *dbDriverStub) ReadAll(collection string) ([]string, error) {
	var resources []string
	for _, c := range d.cards {
		if !strings.HasPrefix(c.Title, collection) {
			b, err := json.Marshal(c)
			if err != nil {
				return resources, err
			}
			resources = append(resources, string(b))
		}
	}
	if len(resources) == 0 {
		return resources, errors.New("collection not found")
	}

	return resources, nil
}

func (d *dbDriverStub) Delete(collection string, resource string) error {
	for i, c := range d.cards {
		g := convertCollectionToGroup(collection)
		cardPath := getCardPath(g, resource)
		if c.Title == cardPath {
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			return nil
		}
	}
	return errors.New("Resource not found")
}

func newRepositoryWithDbAndClockStubs() (*repository, *dbDriverStub) {
	d := &dbDriverStub{}
	c := &clockStub{}
	r := &repository{db: d, clock: c}
	return r, d
}

func newRepositoryWithDbAndClockStubsAndCards() (*repository, *dbDriverStub) {
	r, db := newRepositoryWithDbAndClockStubs()
	db.cards = []Card{
		{Title: "Subject1", Desc: "Value1"},
		{Title: "Subject2", Desc: "Value2"},
		{Title: "Group.Subject1", Desc: "Value1"},
		{Title: "Group.Subject2", Desc: "Value2"},
		{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
		{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
	}
	return r, db
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
			r, db := newRepositoryWithDbAndClockStubs()
			err := r.AddCard(tt.group, tt.card)

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
			r, db := newRepositoryWithDbAndClockStubs()
			var err error
			for _, c := range tt.cards {
				err = r.AddCard(tt.group, c)
			}

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
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
		group   string
		card    deleting.Card
		want    []Card
		wantErr error
	}{
		{
			name:  "Normal",
			group: "Group",
			card:  deleting.Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Card not found",
			group: "Group",
			card:  deleting.Card{Title: "Subject3"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: ErrCardNotFound,
		},
		{
			name:  "Empty Group",
			group: "",
			card:  deleting.Card{Title: "Subject1"},
			want: []Card{
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := newRepositoryWithDbAndClockStubsAndCards()
			err := r.DeleteCard(tt.group, tt.card)
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
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
		},
		{
			name:  "Sub Group",
			group: "Group.SubGroup",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
			},
		},
		{
			name:  "One Card Not Found",
			group: "Group",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
				{Title: "Subject3"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
		},
		{
			name:  "Empty Group",
			group: "",
			cards: []deleting.Card{
				{Title: "Subject1"},
				{Title: "Subject2"},
			},
			want: []Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := newRepositoryWithDbAndClockStubsAndCards()

			for _, c := range tt.cards {
				r.DeleteCard(tt.group, c)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
			name:  "Normal",
			group: "Group",
			want: []getting.Card{
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Sub Group",
			group: "Group.SubGroup",
			want: []getting.Card{
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Empty Group",
			group: "",
			want: []getting.Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := newRepositoryWithDbAndClockStubsAndCards()
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
			name:  "Normal",
			group: "Group",
			card:  updating.Card{Title: "Subject1", Desc: "Value2"},
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
			name:  "Sub Group",
			group: "Group.SubGroup",
			card:  updating.Card{Title: "Subject1", Desc: "Value2"},
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
			name:  "Empty Desc",
			group: "Group",
			card:  updating.Card{Title: "Subject1", Desc: ""},
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
			name:  "Card Not Found",
			group: "Group",
			card:  updating.Card{Title: "Subject3", Desc: "Value3"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
			wantErr: ErrCardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := newRepositoryWithDbAndClockStubsAndCards()
			err := r.UpdateCard(tt.group, tt.card)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
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
				{Title: "Subject1", Desc: "Value2"},
				{Title: "Subject2", Desc: "Value3"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value2"},
				{Title: "Group.Subject2", Desc: "Value3"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
		},
		{
			name:  "Sub Group",
			group: "Group.SubGroup",
			cards: []updating.Card{
				{Title: "Subject1", Desc: "Value2"},
				{Title: "Subject2", Desc: "Value3"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value3"},
			},
		},
		{
			name:  "One Card Not Found",
			group: "Group",
			cards: []updating.Card{
				{Title: "Subject1", Desc: "Value2"},
				{Title: "Subject3", Desc: "Value4"},
			},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Subject2", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value2"},
				{Title: "Group.Subject2", Desc: "Value2"},
				{Title: "Group.SubGroup.Subject1", Desc: "Value1"},
				{Title: "Group.SubGroup.Subject2", Desc: "Value2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, db := newRepositoryWithDbAndClockStubsAndCards()

			for _, c := range tt.cards {
				r.UpdateCard(tt.group, c)
			}

			if !reflect.DeepEqual(tt.want, db.cards) {
				t.Errorf("Incorrect cards. Want %v, got %v", tt.want, db.cards)
			}
		})
	}
}
