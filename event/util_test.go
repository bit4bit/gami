package event

import (
	"reflect"
	"testing"
)

//Util para hacer pruebas de eventos
func testEvent(t *testing.T, fixture map[string]string, evtype interface{}) {
	value := reflect.ValueOf(evtype)
	for _, v := range fixture {
		field := value.FieldByName(v)

		if !field.IsValid() || field.String() == "" {
			t.Log("Not Cast Field:", v)
			t.Fail()
		}
	}
}
