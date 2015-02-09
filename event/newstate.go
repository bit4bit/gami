// Package event for AMI
package event

// Newstate triggered when a channel changes its status.
type Newstate struct {
	Privilege         []string
	Channel           string `AMI:"Channel"`
	ChannelState      string `AMI:"Channelstate"`
	ChannelStateDesc  string `AMI:"Channelstatedesc"`
	CallerIDNum       string `AMI:"Calleridnum"`
	CallerIDName      string `AMI:"Calleridname"`
	UniqueID          string `AMI:"Uniqueid"`
	ConnectedLineNum  string `AMI:"Connectedlinenum"`
	ConnectedLineName string `AMI:"Connectedlinename"`
}

func init() {
	eventTrap["Newstate"] = Newstate{}
}
