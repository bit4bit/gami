//Event decoder
//This Build Type of Event received
package event


import (
	"github.com/bit4bit/GAMI"
	"fmt"
	"reflect"
	"strconv"
)

//used internal for trap events and cast
var eventTrap = make(map[string]interface{})

//Build a new event Type
//if not return AMIEvent
func New(event *gami.AMIEvent) interface{} {
	if intf, ok := eventTrap[event.Id]; ok {
		return  build(event, &intf)
	}
	return *event
}

func build(event *gami.AMIEvent, klass *interface{}) interface{} {

	typ := reflect.TypeOf(*klass)
	value := reflect.ValueOf(*klass)
	ret := reflect.New(typ).Elem()
	for ix := 0; ix < value.NumField(); ix++ {
		field := ret.Field(ix)
		tfield := typ.Field(ix)

		if tfield.Name == "Privilege" {
			field.Set(reflect.ValueOf(event.Privilege))
			continue;
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(event.Params[tfield.Tag.Get("AMI")])
		case reflect.Int64:
			vint,  _ := strconv.Atoi(event.Params[tfield.Tag.Get("AMI")])
			field.SetInt(int64(vint))
		default:
			fmt.Print(ix, tfield.Tag.Get("AMI"), ":", field, "\n")			
		}

	}
	return ret.Interface()
}
