// Package event for AMI
package event

// Newexten triggered when a new extension is accessed.
type Newexten struct {
	Privilege       []string
	Channel         string `AMI:"Channel"`
	Extension       string `AMI:"Extension"`
	Context         string `AMI:"Context"`
	Priority        string `AMI:"Priority"`
	Application     string `AMI:"Application"`
	ApplicationData string `AMI:"Appdata"`
	UniqueID        string `AMI:"Uniqueid"`
}

func init() {
	eventTrap["Newexten"] = Newexten{}
}
