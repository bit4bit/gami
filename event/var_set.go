//Event triggered when a variable is set via agi or dialplan.
package event

type VarSet struct {
	Privilege    []string
	Channel      string `AMI:"Channel"`
	VariableName string `AMI:"Variable"`
	Value        string `AMI:"Value"`
	UniqueID     string `AMI:"Uniqueid"`
}

func init() {
	eventTrap["VarSet"] = VarSet{}
}
