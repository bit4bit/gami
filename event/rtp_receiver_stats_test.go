package event

import (
	"github.com/bit4bit/gami"
	"testing"
)

func TestRTPReceiverStats(t *testing.T) {
	fixture := map[string]string{
		"Ssrc":            "SSRC",
		"Receivedpackets": "ReceivedPackets",
		"Lostpackets":     "LostPackets",
		"Jitter":          "Jitter",
		"Transit":         "Transit",
		"Rrcount":         "RRCount",
	}

	ev := gami.AMIEvent{
		Id:        "RTPReceiverStats",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(RTPReceiverStats); !ok {
		t.Fatal("PeerStatus type assertion")
	}

	testEvent(t, fixture, evtype)
}
