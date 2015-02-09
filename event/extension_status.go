// Package event for AMI
package event

// ExtensionStatus triggered when an extension changes its status.
type ExtensionStatus struct {
	Privilege []string
	Extension string `AMI:"Exten"`
	Context   string `AMI:"Context"`
	Hint      string `AMI:"Hint"`
	Status    string `AMI:"Status"`
}

func init() {
	eventTrap["ExtensionStatus"] = ExtensionStatus{}
}
