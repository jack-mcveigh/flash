package json

import "time"

type Card struct {
	Title   string
	Desc    string
	Created time.Time
	Updated time.Time
}
