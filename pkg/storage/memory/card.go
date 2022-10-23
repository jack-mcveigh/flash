package memory

import "time"

type Card struct {
	Title   string
	Desc    string
	Created time.Time
	Updated time.Time
}
