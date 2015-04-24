package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestHangupEvent(t *testing.T) {
	fixture := map[string]string{
		"Channel":      "Channel",
		"Calleridnum":  "CallerIDNum",
		"Calleridname": "CallerIDName",
		"Uniqueid":     "UniqueID",
		"Cause":        "Cause",
		"Cause-Text":   "CauseText",
	}
	ev := gami.AMIEvent{
		ID:        "Hangup",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Hangup); !ok {
		t.Log("Hangup type assertion")
		t.Fail()
	}

	testEvent(t, fixture, evtype)
}
