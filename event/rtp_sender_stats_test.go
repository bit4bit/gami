package event

import (
	"github.com/bit4bit/gami"
	"testing"
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
		Id:        "RTPSenderStats",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(RTPSenderStats); !ok {
		t.Fatal("PeerStatus type assertion")
	}

	testEvent(t, fixture, evtype)
}
