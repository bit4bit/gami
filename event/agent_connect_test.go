package event

import (
	"github.com/bit4bit/gami"
	"testing"
)


func TestAgentConnect(t *testing.T) {
	fixture := map[string]string{
		"Holdtime": "HoldTime",
		"Bridgedchannel": "BridgedChannel",
		"Ringtime": "RingTime",
		"Member": "Member",
		"Membername": "MemberName",
		"Queue": "Queue",
		"Uniqueid": "UniqueID",
		"Channel": "Channel",
	}

	ev := gami.AMIEvent{
		Id: "AgentConnect",
		Privilege: []string {"all"},
		Params: fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(AgentConnect); !ok {
		t.Fatal("AgentConnect type assertion")
	}

	testEvent(t, fixture, evtype)
}
