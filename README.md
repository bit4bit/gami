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
	"time"
)

func main() {
	ami, err := gami.Dial("127.0.0.1:5038")
	if err != nil {
		log.Fatal(err)
	}
	ami.Run()
	defer ami.Close()
	
	//install manager
	go func() {
		for {
			select {
			//handle network errors
			case err := <-ami.NetError:
				log.Println("Network Error:", err)
				//try new connection every second
				<-time.After(time.Second)
				if err := ami.Reconnect(); err == nil {
					//call start actions
					ami.Action("Events", gami.Params{"EventMask": "on"})
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
		log.Fatal(err)
	}
	
	
	if _, errPing := ami.Action("Ping", nil); errPing != nil {
		log.Fatal(errPing)
	}
	
	//async actions
	rsPing, rsErr := ami.AsyncAction("Ping", gami.Params{"ActionID": "pingo"})
	if rsErr != nil {
		log.Fatal(rsErr)
	}
						
	if _, err := ami.Action("Events", gami.Params{"EventMask":"on"}); err != nil {
		log.Fatal(err)
	}
	
	log.Println("ping:", <-rsPing)
	

}
```

###TLS SUPPORT
In order to use TLS connection to manager interface you could `Dial` with additional parameters
```go
//without TLS
ami, err := gami.Dial("127.0.0.1:5038")

//if certificate is trusted
ami, err := gami.Dial("127.0.0.1:5039", gami.UseTLS)

//if self signed certificate
ami, err := gami.Dial("127.0.0.1:5039", gami.UseTLS, gami.UnsecureTLS)

//if custom tls configuration
ami, err := gami.Dial("127.0.0.1:5039", gami.UseTLSConfig(&tls.Config{}))
```
**WARNING:**
*Only Asterisk >=1.6 supports TLS connection to AMI and
it needs additional configuration(follow the [Asterisk AMI configuration](http://www.asteriskdocs.org/en/3rd_Edition/asterisk-book-html-chunk/AMI-configuration.html) documentation)*

CURRENT EVENT TYPES
====

The events use documentation and struct from *PAMI*.

use **bit4bit/gami/event.New()** for get this struct from raw event

EVENT ID          | TYPE TEST  
----------------  | ---------- 
*Newchannel*      | YES
*Newexten*        | YES
*Newstate*        | YES 
*Dial*            | YES 
*ExtensionStatus* | YES 
*Hangup*          | YES 
*PeerStatus*      | YES
*PeerEntry*	      | YES
*VarSet*          | YES 
*AgentLogin*      | YES
*Agents*          | YES
*AgentLogoff*     | YES
*AgentConnect*    | YES
*RTPReceiverStats* | YES
*RTPSenderStats* | YES
*Bridge* | YES
