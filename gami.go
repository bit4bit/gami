// Package gami provites primitives for interacting with Asterisk AMI
/*

Basic Usage


	ami, err := gami.Dial("127.0.0.1:5038")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
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
		log.Fatal(err)
	}


	if rs, err = ami.Action("Ping", nil); err == nil {
		log.Fatal(rs)
	}

	//or with can do async
	pingResp, pingErr := ami.AsyncAction("Ping", gami.Params{"ActionID": "miping"})
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	if rs, err = ami.Action("Events", ami.Params{"EventMask":"on"}); err != nil {
		fmt.Print(err)
	}

	log.Println("future ping:", <-pingResp)


*/
package gami

import (
	"errors"
	"io"
	"net"
	"net/textproto"
	"strings"
	"syscall"
)

const responseChanGamiID = "gamigeneral"

var errNotEvent = errors.New("Not Event")

// Raise when not response expected protocol AMI
var ErrNotAMI = errors.New("Server not AMI interface")

// Params for the actions
type Params map[string]string

// AMIClient a connection to AMI server
type AMIClient struct {
	conn    *textproto.Conn
	connRaw io.ReadWriteCloser

	address string
	amiUser string
	amiPass string

	// network wait for a new connection
	waitNewConnection chan struct{}

	response map[string]chan *AMIResponse

	// Events for client parse
	Events chan *AMIEvent

	// Error Raise on logic
	Error chan error

	//NetError a network error
	NetError chan error
}

// AMIResponse from action
type AMIResponse struct {
	ID     string
	Status string
	Params map[string]string
}

// AMIEvent it's a representation of Event readed
type AMIEvent struct {
	//Identification of event Event: xxxx
	ID string

	Privilege []string

	// Params  of arguments received
	Params map[string]string
}

// Login authenticate to AMI
func (client *AMIClient) Login(username, password string) error {
	response, err := client.Action("Login", Params{"Username": username, "Secret": password})
	if err != nil {
		return err
	}

	if (*response).Status == "Error" {
		return errors.New((*response).Params["Message"])
	}

	client.amiUser = username
	client.amiPass = password
	return nil
}

// Reconnect the session, autologin if a new network error it put on client.NetError
func (client *AMIClient) Reconnect() error {
	client.conn.Close()
	reconnect, err := Dial(client.address)

	if err != nil {
		client.NetError <- err
		return err
	}

	//new connection
	client.conn = reconnect.conn
	client.connRaw = reconnect.connRaw
	client.waitNewConnection <- struct{}{}

	if err := client.Login(client.amiUser, client.amiPass); err != nil {
		return err
	}

	return nil
}

// AsyncAction return chan for wait response of action with parameter *ActionID* this can be helpful for
// massive actions,
func (client *AMIClient) AsyncAction(action string, params Params) (<-chan *AMIResponse, error) {
	if err := client.conn.PrintfLine("Action: %s", strings.TrimSpace(action)); err != nil {
		return nil, err
	}

	if _, ok := params["ActionID"]; ok {
		params["ActionID"] = responseChanGamiID
	}

	if _, ok := client.response[params["ActionID"]]; !ok {
		client.response[params["ActionID"]] = make(chan *AMIResponse, 1)
	}

	for k, v := range params {
		if err := client.conn.PrintfLine("%s: %s", k, strings.TrimSpace(v)); err != nil {
			return nil, err
		}
	}

	if err := client.conn.PrintfLine(""); err != nil {
		return nil, err
	}

	return client.response[params["ActionID"]], nil
}

// Action send with params
func (client *AMIClient) Action(action string, params Params) (*AMIResponse, error) {
	resp, err := client.AsyncAction(action, params)
	if err != nil {
		return nil, err
	}
	response := <-resp

	return response, nil
}

// Run process socket waiting events and responses
func (client *AMIClient) Run() {
	go func() {
		for {
			data, err := client.conn.ReadMIMEHeader()

			if err != nil {
				switch err {
				case syscall.ECONNABORTED:
					fallthrough
				case syscall.ECONNRESET:
					fallthrough
				case syscall.ECONNREFUSED:
					fallthrough
				case io.EOF:
					client.NetError <- err
					<-client.waitNewConnection
				default:
					client.Error <- err
				}
				continue
			}

			if ev, err := newEvent(&data); err != nil {
				if err != errNotEvent {
					client.Error <- err
				}
			} else {
				client.Events <- ev
			}

			if response, err := newResponse(&data); err == nil {
				client.response[response.ID] <- response
			}

		}
	}()
}

// Close the connection to AMI
func (client *AMIClient) Close() {
	client.Action("Logoff", nil)
	(client.connRaw).Close()
}

//newResponse build a response for action
func newResponse(data *textproto.MIMEHeader) (*AMIResponse, error) {
	if data.Get("Response") == "" {
		return nil, errors.New("Not Response")
	}
	response := &AMIResponse{"", "", make(map[string]string)}
	for k, v := range *data {
		if k == "Response" {
			continue
		}
		response.Params[k] = v[0]
	}
	response.ID = data.Get("Actionid")
	response.Status = data.Get("Response")
	return response, nil
}

//newEvent build event
func newEvent(data *textproto.MIMEHeader) (*AMIEvent, error) {
	if data.Get("Event") == "" {
		return nil, errNotEvent
	}
	ev := &AMIEvent{data.Get("Event"), strings.Split(data.Get("Privilege"), ","), make(map[string]string)}
	for k, v := range *data {
		if k == "Event" || k == "Privilege" {
			continue
		}
		ev.Params[k] = v[0]
	}
	return ev, nil
}

// Dial create a new connection to AMI
func Dial(address string) (*AMIClient, error) {
	connRaw, err := net.Dial("tcp", address)

	if err != nil {
		return nil, err
	}
	conn := textproto.NewConn(connRaw)
	label, err := conn.ReadLine()
	if err != nil {
		return nil, err
	}

	if strings.Contains(label, "Asterisk Call Manager") != true {
		return nil, ErrNotAMI
	}

	client := &AMIClient{
		conn:              conn,
		connRaw:           connRaw,
		address:           address,
		amiUser:           "",
		amiPass:           "",
		waitNewConnection: make(chan struct{}),
		response:          make(map[string]chan *AMIResponse),
		Events:            make(chan *AMIEvent, 100),
		Error:             make(chan error, 1),
		NetError:          make(chan error, 1),
	}
	return client, nil
}

// NewConn create a new connection to AMI
func NewConn(connRaw io.ReadWriteCloser, address string) (*AMIClient, error) {
	conn := textproto.NewConn(connRaw)
	label, err := conn.ReadLine()
	if err != nil {
		return nil, err
	}

	if strings.Contains(label, "Asterisk Call Manager") != true {
		return nil, ErrNotAMI
	}

	client := &AMIClient{
		conn:              conn,
		connRaw:           connRaw,
		address:           address,
		amiUser:           "",
		amiPass:           "",
		waitNewConnection: make(chan struct{}),
		response:          make(map[string]chan *AMIResponse),
		Events:            make(chan *AMIEvent, 100),
		Error:             make(chan error, 1),
		NetError:          make(chan error, 1),
	}
	return client, nil
}
