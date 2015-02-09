// Package event for AMI
package event

// Newchannel triggered when a new channel is created.
type Newchannel struct {
	Privilege        []string
	Channel          string `AMI:"Channel"`
	ChannelState     string `AMI:"Channelstate"`
	ChannelStateDesc string `AMI:"Channelstatedesc"`
	CallerIDNum      string `AMI:"Calleridnum"`
	CallerIDName     string `AMI:"Calleridname"`
	AccountCode      string `AMI:"Accountcode"`
	UniqueID         string `AMI:"Uniqueid"`
	Context          string `AMI:"Context"`
	Extension        string `AMI:"Exten"`
}

func init() {
	eventTrap["Newchannel"] = Newchannel{}
}
