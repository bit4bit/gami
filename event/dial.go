// Package event for AMI
package event

// Dial triggered when a dial is executed.
type Dial struct {
	Privilege    []string
	SubEvent     string `AMI:"Subevent"`
	Channel      string `AMI:"Channel"`
	Destination  string `AMI:"Destination"`
	CallerIDNum  string `AMI:"Calleridnum"`
	CallerIDName string `AMI:"Calleridname"`
	UniqueID     string `AMI:"Uniqueid"`
	DestUniqueID string `AMI:"Destuniqueid"`
	DialString   string `AMI:"Dialstring"`
	DialStatus   string `AMI:"Dialstatus"`
}

func init() {
	eventTrap["Dial"] = Dial{}
}
