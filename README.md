GAMI
====

GO - Asterisk AMI Interface

communicate with the  Asterisk AMI, Actions and Events.

Example connecting to Asterisk and Send Action get Events.

```go
package main
import (
	"log"
	"github.com/bit4bit/gami"
	"github.com/bit4bit/gami/event"
)

func main() {
	ami, err := gami.Dial("127.0.0.1:5038")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	
	//install manager
	go func() {
		for {
			select {
			//handle network errors
			case err := <-ami.NetError:
				log.Println("Network Error:", err)
				//try new connection every second
				for {
					<-time.After(time.Second)
					if err := ami.Reconnect(); err == nil {
						//call start actions
						if _, err := ami.Action("Events", gami.Params{"EventMask": "on"}); err != nil {
							break
						}
					}
				}
				
			case err := <-ami.Error:
				log.Println("error:", err)
			//wait events and process
			case ev := <-ami.Events:
				log.Println("Event Detect: %v", *ev)
				//if want type of events
				log.Println("EventType:", event.New(ev))
			}
		}
	}()
	
	if err := ami.Login("admin", "root"); err != nil {
		fmt.Print(err)
	}
	
	
	if rs, err = ami.Action("Ping", nil); err != nil {
		fmt.Print(rs)
	}
	if rs, err = ami.Action("Events", ami.Params{"EventMask":"on"}); err != nil {
		fmt.Print(err)
	}
	
	ami.Close()
}
```



CURRENT EVENT TYPES
====

The events use documentation and struct from *PAMI*.

EVENT ID          | TYPE TEST  
----------------  | ---------- 
*Newchannel*      | YES
*Newexten*        | YES
*Newstate*        | YES 
*Dial*            | YES 
*ExtensionStatus* | YES 
*Hangup*          | YES 
*PeerStatus*      | YES
*VarSet*          | YES 
*AgentLogin*      | YES
*Agents*          | YES
*AgentLogoff*     | YES
*AgentConnect*    | YES
