package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestAgentConnect(t *testing.T) {
	fixture := map[string]string{
		"Holdtime":       "HoldTime",
		"Bridgedchannel": "BridgedChannel",
		"Ringtime":       "RingTime",
		"Member":         "Member",
		"Membername":     "MemberName",
		"Queue":          "Queue",
		"Uniqueid":       "UniqueID",
		"Channel":        "Channel",
	}

	ev := gami.AMIEvent{
		ID:        "AgentConnect",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(AgentConnect); !ok {
		t.Fatal("AgentConnect type assertion")
	}

	testEvent(t, fixture, evtype)
}
