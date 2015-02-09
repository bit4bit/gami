// Package event for AMI
package event

// Hangup triggered when a hangup is detected.
type Hangup struct {
	Privilege    []string
	Channel      string `AMI:"Channel"`
	CallerIDNum  string `AMI:"Calleridnum"`
	CallerIDName string `AMI:"Calleridname"`
	UniqueID     string `AMI:"Uniqueid"`
	Cause        string `AMI:"Cause"`
	CauseText    string `AMI:"Cause-Text"`
}

func init() {
	eventTrap["Hangup"] = Hangup{}
}
