// Package event for AMI
package event

//Agents trigger for agents
type Agents struct {
	Privilege        []string
	Status           string `AMI:"Status"`
	Agent            string `AMI:"Agent"`
	Name             string `AMI:"Name"`
	Channel          string `AMI:"Channel"`
	LoggedInTime     string `AMI:"Loggedintime"`
	TalkingTo        string `AMI:"Talkingto"`
	TalkingToChannel string `AMI:"Talkingtochannel"`
}

func init() {
	eventTrap["Agents"] = Agents{}
}
