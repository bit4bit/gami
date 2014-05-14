//GAMI
//Library for interacting with Asterisk AMI
package gami

import (
	"errors"
	"net/textproto"
	"strings"
)

var errInvalidLogin error = errors.New("InvalidLogin AMI Interface")
var errBadResponse error = errors.New("Bad Response for action")
var errNotAMI error = errors.New("Server not AMI interface")

type AMIClient struct {
	conn *textproto.Conn

	response chan *AMIResponse

	Events chan *AMIEvent
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

	if _, err := client.Action("Login", map[string]string{"Username": username, "Secret": password}); err != nil {
		return err
	}

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

	return <-client.response, nil
}

//Process socket waiting events and responses
func (client *AMIClient) run() {

	go func() {

		for {

			data, _ := client.conn.ReadMIMEHeader()

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
	conn, err := textproto.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	label, _ := conn.ReadLine()
	if strings.Contains(label, "Asterisk Call Manager") != true {
		return nil, errNotAMI
	}

	client := &AMIClient{conn, make(chan *AMIResponse), make(chan *AMIEvent, 100)}
	client.run()
	return client, nil
}
