//Event Trigger for agents
package event

type Agents struct {
	Privilege []string
	Status string `AMI:"Status"`
	Agent string `AMI:"Agent"`
	Name string `AMI:"Name"`
	Channel string `AMI:"Channel"`
	LoggedInTime string `AMI:"Channel"`
	TalkingTo string `AMI:"TalkingTo"`
	TalkingToChannel string `AMI:"TalkingToChannel"`
}

func init() {
	eventTrap["Agents"] = Agents{}
}
