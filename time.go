package esi

import "time"

// Timestamp is a time.Time
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}
