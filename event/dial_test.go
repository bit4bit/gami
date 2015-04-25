package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestDialEvent(t *testing.T) {
	fixture := map[string]string{
		"Subevent":     "SubEvent",
		"Channel":      "Channel",
		"Destination":  "Destination",
		"Calleridnum":  "CallerIDNum",
		"Calleridname": "CallerIDName",
		"Uniqueid":     "UniqueID",
		"Destuniqueid": "DestUniqueID",
		"Dialstring":   "DialString",
		"Dialstatus":   "DialStatus",
	}

	ev := gami.AMIEvent{
		ID:        "Dial",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Dial); !ok {
		t.Log("Dial type assertion")
		t.Fail()
	}
	testEvent(t, fixture, evtype)

}
