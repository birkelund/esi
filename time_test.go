package esi

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestTimestamp_MarshalJSON(t *testing.T) {
	ts := Timestamp{time.Now()}
	expected := fmt.Sprintf(`"%s"`, ts.Format(time.RFC3339))
	buf, _ := json.Marshal(ts)
	if string(buf) != expected {
		t.Fatalf("expected %q; got %q", expected, string(buf))
	}
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	data := fmt.Sprintf(`"%s"`, ts.Format(time.RFC3339))
	var unmarshalled Timestamp
	json.Unmarshal([]byte(data), &unmarshalled)

	if !unmarshalled.Time.Equal(ts) {
		t.Fatalf("want %q; got %q", ts, unmarshalled)
	}
}

func TestTimestamp_UnmarshalJSON_invalidJSON(t *testing.T) {
	var ts Timestamp
	err := json.Unmarshal([]byte("invalid json"), &ts)
	if err == nil {
		t.Fatal("expected error")
	}

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatal("expected *json.SyntaxError")
	}
}

func TestTimestamp_UnmarshalJSON_invalidTimestampFormat(t *testing.T) {
	var ts Timestamp
	err := json.Unmarshal([]byte(`"Mon, 02 Jan 2006 15:04:05 MST"`), &ts)
	if err == nil {
		t.Fatal("expected error")
	}

	if _, ok := err.(*time.ParseError); !ok {
		t.Fatal("expected *time.ParseError")
	}
}
