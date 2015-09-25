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
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/textproto"
	"strings"
	"sync"
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
	conn             *textproto.Conn
	connRaw          io.ReadWriteCloser
	mutexAsyncAction *sync.RWMutex

	address     string
	amiUser     string
	amiPass     string
	useTLS      bool
	unsecureTLS bool

	// TLSConfig for secure connections
	tlsConfig *tls.Config

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

func UseTLS(c *AMIClient) {
	c.useTLS = true
}

func UseTLSConfig(config *tls.Config) func(*AMIClient) {
	return func(c *AMIClient) {
		c.tlsConfig = config
		c.useTLS = true
	}
}

func UnsecureTLS(c *AMIClient) {
	c.unsecureTLS = true
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
	err := client.NewConn()

	if err != nil {
		client.NetError <- err
		return err
	}

	client.waitNewConnection <- struct{}{}

	if err := client.Login(client.amiUser, client.amiPass); err != nil {
		return err
	}

	return nil
}

// AsyncAction return chan for wait response of action with parameter *ActionID* this can be helpful for
// massive actions,
func (client *AMIClient) AsyncAction(action string, params Params) (<-chan *AMIResponse, error) {
	var output string
	client.mutexAsyncAction.Lock()
	defer client.mutexAsyncAction.Unlock()

	output = fmt.Sprintf("Action: %s\r\n", strings.TrimSpace(action))
	if params == nil {
		params = Params{}
	}
	if _, ok := params["ActionID"]; !ok {
		params["ActionID"] = responseChanGamiID + "_" + fmt.Sprintf("%d", rand.Intn(1000000))
	}

	if _, ok := client.response[params["ActionID"]]; !ok {
		client.response[params["ActionID"]] = make(chan *AMIResponse, 1)
	}
	for k, v := range params {
		output = output + fmt.Sprintf("%s: %s\r\n", k, strings.TrimSpace(v))
	}
	if err := client.conn.PrintfLine("%s", output); err != nil {
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
				client.notifyResponse(response)
			}

		}
	}()
}

// Close the connection to AMI
func (client *AMIClient) Close() {
	client.Action("Logoff", nil)
	(client.connRaw).Close()
}

func (client *AMIClient) notifyResponse(response *AMIResponse) {
	go func() {
		client.mutexAsyncAction.RLock()
		client.response[response.ID] <- response
		close(client.response[response.ID])
		delete(client.response, response.ID)
		client.mutexAsyncAction.RUnlock()
	}()
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
func Dial(address string, options ...func(*AMIClient)) (*AMIClient, error) {
	client := &AMIClient{
		address:           address,
		amiUser:           "",
		amiPass:           "",
		mutexAsyncAction:  new(sync.RWMutex),
		waitNewConnection: make(chan struct{}),
		response:          make(map[string]chan *AMIResponse),
		Events:            make(chan *AMIEvent, 100),
		Error:             make(chan error, 1),
		NetError:          make(chan error, 1),
		useTLS:            false,
		unsecureTLS:       false,
		tlsConfig:         new(tls.Config),
	}
	for _, op := range options {
		op(client)
	}
	err := client.NewConn()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewConn create a new connection to AMI
func (client *AMIClient) NewConn() (err error) {
	if client.useTLS {
		client.tlsConfig.InsecureSkipVerify = client.unsecureTLS
		client.connRaw, err = tls.Dial("tcp", client.address, client.tlsConfig)
	} else {
		client.connRaw, err = net.Dial("tcp", client.address)
	}

	if err != nil {
		return err
	}

	client.conn = textproto.NewConn(client.connRaw)
	label, err := client.conn.ReadLine()
	if err != nil {
		return err
	}

	if strings.Contains(label, "Asterisk Call Manager") != true {
		return ErrNotAMI
	}

	return nil
}
