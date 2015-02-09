// Package event for AMI
package event

// AgentConnect triggered when an agent connects.
type AgentConnect struct {
	Privilege      []string
	HoldTime       string `AMI:"Holdtime"`
	BridgedChannel string `AMI:"Bridgedchannel"`
	RingTime       string `AMI:"Ringtime"`
	Member         string `AMI:"Member"`
	MemberName     string `AMI:"Membername"`
	Queue          string `AMI:"Queue"`
	UniqueID       string `AMI:"Uniqueid"`
	Channel        string `AMI:"Channel"`
}

func init() {
	eventTrap["AgentConnect"] = AgentConnect{}
}
