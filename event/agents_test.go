package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestAgentsEvent(t *testing.T) {
	fixture := map[string]string{
		"Status":           "Status",
		"Agent":            "Agent",
		"Name":             "Name",
		"Channel":          "Channel",
		"Loggedintime":     "LoggedInTime",
		"Talkingto":        "TalkingTo",
		"Talkingtochannel": "TalkingToChannel",
	}

	ev := gami.AMIEvent{
		ID:        "Agents",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Agents); !ok {
		t.Fatal("Agents type assertion")
	}
	testEvent(t, fixture, evtype)
}
