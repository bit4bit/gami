//Event triggered when a dial is executed.
package event

type Dial struct {
	Privilege    []string
	SubEvent     string `AMI:"Subevent"`
	Channel      string `AMI:"Channel"`
	Destination  string `AMI:"Destination"`
	CallerIDNum  string `AMI:"Calleridnum"`
	CallerIDName string `AMI:"Calleridname"`
	UniqueID     string `AMI:"UniqueID"`
	DestUniqueID string `AMI:"Destuniqueid"`
	DialString   string `AMI:"Dialstring"`
	DialStatus   string `AMI:"Dialstatus"`
}

func init() {
	eventTrap["Dial"] = Dial{}
}
