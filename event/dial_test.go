package event

import (
	"github.com/bit4bit/GAMI"
	"testing"
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
		Id:        "Dial",
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
