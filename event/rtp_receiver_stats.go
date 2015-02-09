// Package event for AMI
package event

// RTPReceiverStats triggered when exchanging rtp stats.
type RTPReceiverStats struct {
	Privilege       []string
	SSRC            string `AMI:"Ssrc"`
	ReceivedPackets int64  `AMI:"Receivedpackets"`
	LostPackets     int64  `AMI:"Lostpackets"`
	Jitter          string `AMI:"Jitter"`
	Transit         string `AMI:"Transit"`
	RRCount         string `AMI:"Rrcount"`
}

func init() {
	eventTrap["RTPReceiverStats"] = RTPReceiverStats{}
}
