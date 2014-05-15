package event

import (
	"github.com/bit4bit/GAMI"
	"testing"
)

func TestNewstateEvent(t *testing.T) {
	fixture := map[string]string{
		"Channel":           "Channel",
		"Channelstate":      "ChannelState",
		"Channelstatedesc":  "ChannelStateDesc",
		"Calleridnum":       "CallerIDNum",
		"Calleridname":      "CallerIDName",
		"Uniqueid":          "UniqueID",
		"Connectedlinenum":  "ConnectedLineNum",
		"Connectedlinename": "ConnectedLineName",
	}

	ev := gami.AMIEvent{
		Id:        "Newstate",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Newstate); !ok {
		t.Log("Newstate type assertion")
		t.Fail()
	}

	testEvent(t, fixture, evtype)
}
