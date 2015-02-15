// Package event for AMI
package event

// PeerEntry triggered for each peer when an action Sippeers is issued.
type PeerEntry struct {
	Privilege         []string
	ChannelType       string `AMI:"Channeltype"`
	ObjectName        string `AMI:"Objectname"`
	ChannelObjectType string `AMI:"Chanobjecttype"`
	IPAddress         string `AMI:"Ipaddress"`
	IPPort            string `AMI:"Ipport"`
	Dynamic           string `AMI:"Dynamic"`
	NatSupport        string `AMI:"Natsupport"`
	VideoSupport      string `AMI:"Videosupport"`
	TextSupport       string `AMI:"Textsupport"`
	ACL               string `AMI:"Acl"`
	Status            string `AMI:"Status"`
	RealtimeDevice    string `AMI:"Realtimedevice"`
}

func init() {
	//Register ID Event for cast when detect
	eventTrap["PeerEntry"] = PeerEntry{}
}
