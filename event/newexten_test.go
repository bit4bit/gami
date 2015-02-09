package event

import (
	"github.com/bit4bit/gami"
	"testing"
)

func TestNewexten(t *testing.T) {
	fixture := map[string]string{
		"Channel":     "Channel",
		"Extension":   "Extension",
		"Context":     "Context",
		"Priority":    "Priority",
		"Application": "Application",
	}

	ev := gami.AMIEvent{
		Id:        "Newexten",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Newexten); !ok {
		t.Fatal("Newexten type assertion")
	}

	testEvent(t, fixture, evtype)
}
