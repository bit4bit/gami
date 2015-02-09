package event

import (
	"github.com/bit4bit/gami"
	"testing"
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
		Id:        "Newchannel",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(Newchannel); !ok {
		t.Fatal("Newchannel type assertion")
	}

	testEvent(t, fixture, evtype)
}
