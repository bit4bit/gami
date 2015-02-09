// Package event for AMI
package event

// AgentLogoff triggered when an agent logs off.
type AgentLogoff struct {
	Privilege []string
	Agent     string `AMI:"Agent"`
	UniqueID  string `AMI:"Uniqueid"`
	LoginTime string `AMI:"Logintime"`
}

func init() {
	eventTrap["AgentLogoff"] = AgentLogoff{}
}
