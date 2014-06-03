//Event triggered when an agent logs off.
package event

type AgentLogoff struct {
	Privilege []string
	Agent string `AMI:"Agent"`
	UniqueID string `AMI:"Uniqueid"`
	LoginTime string `AMI:"Logintime"`
}

func init() {
	eventTrap["AgentLogoff"] = AgentLogoff{}
}
