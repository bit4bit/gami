package event

import (
	"github.com/bit4bit/GAMI"
	"testing"
)

func TestVarSetEvent(t *testing.T) {
	fixture := map[string]string{
		"Channel" : "Channel",
		"Variable" : "VariableName",
		"Value" : "Value",
		"Uniqueid" : "UniqueID",
	}

	ev := gami.AMIEvent{
		Id: "VarSet",
		Privilege: []string{"all"},
		Params: fixture,
	}
	
	evtype := New(&ev)
	if _, ok := evtype.(VarSet); !ok {
		t.Log("VarSet type assertion")
		t.Fail()
	}

	testEvent(t, fixture, evtype)
}
