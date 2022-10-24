package updating

import (
	"reflect"
	"testing"
)

type repositoryStub struct {
	cards []Card
}

func (r *repositoryStub) UpdateCard(g string, c Card) error {
	if g != "" {
		c.Title = g + "." + c.Title
	}

	for i := range r.cards {
		if r.cards[i].Title == c.Title {
			r.cards[i].Desc = c.Desc
			return nil
		}
	}
	return ErrCardNotFound
}

func TestUpdateCard(t *testing.T) {
	tests := []struct {
		name    string
		group   string
		card    Card
		want    []Card
		wantErr error
	}{
		{
			name:  "Normal",
			group: "Group",
			card:  Card{Title: "Subject1", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value2"},
			},
			wantErr: nil,
		},
		{
			name:  "Empty Desc",
			group: "Group",
			card:  Card{Title: "Subject1", Desc: ""},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: ""},
			},
			wantErr: nil,
		},

		{
			name:  "No Group",
			group: "",
			card:  Card{Title: "Subject1", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value2"},
				{Title: "Group.Subject1", Desc: "Value1"},
			},
			wantErr: nil,
		},
		{
			name:  "Card Not Found",
			group: "Group",
			card:  Card{Title: "Subject2", Desc: "Value2"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
			},
			wantErr: ErrCardNotFound,
		},
		{
			name:  "Empty Title",
			group: "Group",
			card:  Card{Title: "", Desc: "Value"},
			want: []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
			},
			wantErr: ErrCardEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repositoryStub{}
			repo.cards = []Card{
				{Title: "Subject1", Desc: "Value1"},
				{Title: "Group.Subject1", Desc: "Value1"},
			}
			us := New(repo)
			err := us.UpdateCard(tt.group, tt.card)

			if err != tt.wantErr {
				t.Errorf("Incorrect error. Want %v, got %v", tt.wantErr, err)
			}

			if !reflect.DeepEqual(tt.want, repo.cards) {
				t.Errorf("Incorrect repo.cards. Want %v, got %v", tt.want, repo.cards)
			}
		})
	}
}
