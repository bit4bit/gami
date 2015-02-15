// Package event for AMI
package event

type Bridge struct {
	Privilege   []string
	BridgeState string `AMI:"Bridgestate"`
	BridgeType  string `AMI:"Bridgetype"`
	Channel1    string `AMI:"Channel1"`
	Channel2    string `AMI:"Channel2"`
	CallerID1   string `AMI:"Callerid1"`
	CallerID2   string `AMI:"Callerid2"`
	UniqueID1   string `AMI:"Uniqueid1"`
	UniqueID2   string `AMI:"Uniqueid2"`
}

func init() {
	eventTrap["Bridge"] = Bridge{}
}
