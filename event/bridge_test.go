// Package event for AMI
package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestBridge(t *testing.T) {
	fixture := map[string]string{
		"Bridgestate": "BridgeState",
		"Bridgetype":  "BridgeType",
		"Channel1":    "Channel1",
		"Channel2":    "Channel2",
		"Callerid1":   "CallerID1",
		"Callerid2":   "CallerID2",
		"Uniqueid1":   "UniqueID1",
		"Uniqueid2":   "UniqueID2",
	}

	ev := gami.AMIEvent{
		ID:        "Bridge",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Bridge); !ok {
		t.Fatal("Bridge type assertion")
	}

	testEvent(t, fixture, evtype)
}
