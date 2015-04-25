package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestNewchannel(t *testing.T) {
	fixture := map[string]string{
		"Channel":          "Channel",
		"Channelstate":     "ChannelState",
		"Channelstatedesc": "ChannelStateDesc",
		"Calleridnum":      "CallerIDNum",
		"Calleridname":     "CallerIDName",
		"Accountcode":      "AccountCode",
		"Uniqueid":         "UniqueID",
		"Context":          "Context",
		"Exten":            "Extension",
	}

	ev := gami.AMIEvent{
		ID:        "Newchannel",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Newchannel); !ok {
		t.Fatal("Newchannel type assertion")
	}

	testEvent(t, fixture, evtype)
}
