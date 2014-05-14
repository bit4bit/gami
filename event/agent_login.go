//Event trigger when agent logs in
package event

type AgentLogin struct {
	Privilege []string
	Agent     string `AMI:"Agent"`
	UniqueID  string `AMI:"Uniqueid"`
	Channel   string `AMI:"Channel"`
}

func init() {
	eventTrap["AgentLogin"] = AgentLogin{}
}
