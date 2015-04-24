package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestRTPSenderStats(t *testing.T) {
	fixture := map[string]string{
		"Ssrc":        "SSRC",
		"Sendpackets": "SendPackets",
		"Lostpackets": "LostPackets",
		"Jitter":      "Jitter",
		"Rtt":         "RTT",
		"Srcount":     "SRCount",
	}

	ev := gami.AMIEvent{
		ID:        "RTPSenderStats",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(RTPSenderStats); !ok {
		t.Fatal("PeerStatus type assertion")
	}

	testEvent(t, fixture, evtype)
}
