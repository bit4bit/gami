// Package event for AMI
package event

// RTPSenderStats triggered when exchanging rtp stats.
type RTPSenderStats struct {
	Privilege   []string
	SSRC        string `AMI:"Ssrc"`
	SendPackets int64  `AMI:"Sendpackets"`
	LostPackets int64  `AMI:"Lostpackets"`
	Jitter      string `AMI:"Jitter"`
	RTT         string `AMI:"Rtt"`
	SRCount     string `AMI:"Srcount"`
}

func init() {
	eventTrap["RTPSenderStats"] = RTPSenderStats{}
}
