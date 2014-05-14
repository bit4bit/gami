//Triggered when an extension changes its status.
package event

type ExtensionStatus struct {
	Privilege []string
	Extension string `AMI:"Exten"`
	Context string `AMI:"Context"`
	Hint string `AMI:"String"`
	Status string `AMI:"String"`
}

func init() {
	eventTrap["ExtensionStatus"] = ExtensionStatus{}
}
