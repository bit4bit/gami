package event

import (
	"github.com/bit4bit/gami"
	"testing"
)

func TestAgentLogin(t *testing.T) {
	fixture := map[string]string{
		"Agent": "Agent",
		"Uniqueid": "UniqueID",
		"Channel": "Channel",
	}

	ev := gami.AMIEvent{
		Id: "AgentLogin",
		Privilege: []string{"all"},
		Params: fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(AgentLogin); !ok {
		t.Fatal("AgentLogin type assertion")
	}

	testEvent(t, fixture, evtype)
}
