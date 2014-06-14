//GAMI
//Library for interacting with Asterisk AMI
package gami

import (
	"errors"
	"net/textproto"
	"strings"
	"net"
	"time"
)

var errInvalidLogin error = errors.New("InvalidLogin AMI Interface")
var errBadResponse error = errors.New("Bad Response for action")
var errNotAMI error = errors.New("Server not AMI interface")

type AMIClient struct {
	conn *textproto.Conn
	conn_raw *net.Conn

	address string
	amiUser string
	amiPass string
	
	//Timeout on read
	ReadTimeout time.Duration

	response chan *AMIResponse

	//Events for client parse
	Events chan *AMIEvent

	//Raise error
	Error chan error
}

type AMIResponse struct {
	Status string
	Params map[string]string
}

//Representation of Event readed
type AMIEvent struct {
	//Identification of event Event: xxxx
	Id string

	Privilege []string

	//Rest of arguments received
	Params map[string]string
}

func (client *AMIClient) Login(username, password string) error {

	 response, err := client.Action("Login", map[string]string{"Username": username, "Secret": password}); 
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

//Reconnect the session, autologin
func (client *AMIClient) Reconnect() error {
	client.conn.Close()

	reconnect, err := Dial(client.address)
	if err != nil {
		return err
	}
	if err := reconnect.Login(client.amiUser, client.amiPass); err != nil {
		return err
	}

	client.conn = reconnect.conn
	client.conn_raw = reconnect.conn_raw
	return nil
}

//Send a single action
func (client *AMIClient) Action(action string, params map[string]string) (*AMIResponse, error) {


	if err := client.conn.PrintfLine("Action: %s", strings.TrimSpace(action)); err != nil {
		return nil, err
	}

	for k, v := range params {
		if err := client.conn.PrintfLine("%s: %s", k, strings.TrimSpace(v)); err != nil {
			return nil, err
		}
	}

	if err := client.conn.PrintfLine(""); err != nil {
		return nil, err
	}

	response := <-client.response
	return response, nil
}

//Process socket waiting events and responses
func (client *AMIClient) run() {

	go func() {

		for {
			if client.ReadTimeout > 0 {
				//close the connection if not read something
				(*client.conn_raw).SetReadDeadline(time.Now().Add(client.ReadTimeout))			
			}

			data, err := client.conn.ReadMIMEHeader()

			if err != nil {
				client.Error <- err
				continue
			}

			if ev, err := newEvent(&data); err == nil {
				client.Events <- ev
			}
			if response, err := newResponse(&data); err == nil {
				client.response <- response
			}
		}
	}()
}

func (client *AMIClient) Close() {
	client.Action("Logoff", nil)
	client.conn.Close()
}

func newResponse(data *textproto.MIMEHeader) (*AMIResponse, error) {
	if data.Get("Response") == "" {
		return nil, errors.New("Not Response")
	}
	response := &AMIResponse{"", make(map[string]string)}
	for k, v := range *data {
		if k == "Response" {
			continue
		}
		response.Params[k] = v[0]
	}
	response.Status = data.Get("Response")
	return response, nil
}

func newEvent(data *textproto.MIMEHeader) (*AMIEvent, error) {
	if data.Get("Event") == "" {
		return nil, errors.New("Not Event")
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

func Dial(address string) (*AMIClient, error) {
	conn_raw, err := net.Dial("tcp", address)
	
	if err != nil {
		return nil, err
	}
	conn := textproto.NewConn(conn_raw)
	label, err := conn.ReadLine()
	if err != nil {
		return nil, err
	}
	
	if strings.Contains(label, "Asterisk Call Manager") != true {
		return nil, errNotAMI
	}

	client := &AMIClient{
		conn: conn, 
		conn_raw: &conn_raw, 
		address: address,
		amiUser: "",
		amiPass: "",
		response: make(chan *AMIResponse),
		ReadTimeout: 0,
		Events: make(chan *AMIEvent, 100),
		Error: make(chan error)}
	client.run()
	return client, nil
}
