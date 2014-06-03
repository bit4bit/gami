package event

import (
	"reflect"
	"testing"
)

//Util para hacer pruebas de eventos
func testEvent(t *testing.T, fixture map[string]string, evtype interface{}) {
	value := reflect.ValueOf(evtype)
	typ := reflect.TypeOf(evtype)
	for k, v := range fixture {
		field := value.FieldByName(v)
		tfield, _ := typ.FieldByName(v)

		if tfield.Tag.Get("AMI") != k {
			t.Fatal("Not Cast AMI Field:", k, " from", v)
		}

		if !field.IsValid() || field.String() == "" {
			t.Fatal("Not Cast Field:", v)
		}
	}
}
