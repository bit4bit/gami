package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestPeerStatus(t *testing.T) {
	fixture := map[string]string{
		"Channeltype": "ChannelType",
		"Peer":        "Peer",
		"Peerstatus":  "PeerStatus",
	}

	ev := gami.AMIEvent{
		ID:        "PeerStatus",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(PeerStatus); !ok {
		t.Fatal("PeerStatus type assertion")
	}

	testEvent(t, fixture, evtype)
}
