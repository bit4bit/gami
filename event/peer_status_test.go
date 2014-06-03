package event

import (
	"github.com/bit4bit/gami"
	"testing"
)

func TestPeerStatus(t *testing.T) {
	fixture := map[string]string {
		"Channeltype": "ChannelType",
		"Peer": "Peer",
		"Peerstatus": "PeerStatus",
	}

	ev := gami.AMIEvent{
		Id: "PeerStatus",
		Privilege: []string{"all"},
		Params: fixture,
	}
	
	evtype := New(&ev)
	if _, ok := evtype.(PeerStatus); !ok {
		t.Fatal("PeerStatus type assertion")
	}

	testEvent(t, fixture, evtype)
}
