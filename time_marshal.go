package esi

import (
	"fmt"
	"time"
)

// MarshalJSON implements the json.Marshaller interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, time.Time(t.Time).Format(time.RFC3339))

	return []byte(s), nil
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	s := string(b)

	s = s[1 : len(s)-1]

	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = Timestamp{parsed}

	return nil
}
