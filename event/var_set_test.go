package event

import (
	"testing"

	"github.com/bit4bit/gami"
)

func TestVarSetEvent(t *testing.T) {
	fixture := map[string]string{
		"Channel":  "Channel",
		"Variable": "VariableName",
		"Value":    "Value",
		"Uniqueid": "UniqueID",
	}

	ev := gami.AMIEvent{
		ID:        "VarSet",
		Privilege: []string{"all"},
		Params:    fixture,
	}

	evtype := New(&ev)
	if _, ok := evtype.(VarSet); !ok {
		t.Log("VarSet type assertion")
		t.Fail()
	}

	testEvent(t, fixture, evtype)
}
