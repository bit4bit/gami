// Package event for AMI
package event

// PeerStatus trigger when a peers change status
type PeerStatus struct {
	Privilege   []string
	ChannelType string `AMI:"Channeltype"`
	Peer        string `AMI:"Peer"`
	PeerStatus  string `AMI:"Peerstatus"`
}

func init() {
	//Register ID Event for cast when detect
	eventTrap["PeerStatus"] = PeerStatus{}
}
