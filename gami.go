//Implementa servico para la lectura de eventos de Asterisk *AMI*
package gami

import (
	"net/textproto"
	"strings"
	"errors"
	"fmt"
)

var errInvalidLogin error = errors.New("InvalidLogin AMI Interface")
var errBadResponse error = errors.New("Bad Response for action")
var errNotAMI error = errors.New("Server not AMI interface")

type AMIClient struct {
	conn *textproto.Conn
	
	response chan AMIResponse
	
	//Se encolan los eventos 
	//para ser procesados por el cliente de la API
	Events chan AMIEvent
}

type AMIResponse struct {
	Status string
	Params map[string]string
}

type AMIEvent struct {
	Id string
	Privilege []string
	Params map[string]string
}


func (client *AMIClient) Login(username, password string) error  {
	
	if _, err := client.Action("Login", map[string]string{"Username": username, "Secret": password}) ; err != nil {
		return err
	}

	return nil
}

//
//@todo con la Action:Logoff se cierra conexion y quedan bloqueados los
func (client *AMIClient) Action(action string, params map[string]string) (AMIResponse, error) {

	fmt.Print("AMIClient::Action: ", action, "\n")
	
	if err := client.conn.PrintfLine("Action: %s", strings.TrimSpace(action)); err != nil {
		fmt.Printf("ERRROR:%v", err)

	}
	for k,v := range params {
		client.conn.PrintfLine("%s: %s", k, strings.TrimSpace(v))
	}
	client.conn.PrintfLine("")
	
	fmt.Print("AMIClient::Action::END\n")
	response := <-client.response
	return response, nil
}


//Va enviando eventos obtenidos
//!importante debe ser llamada despues de *Login*
//para obtener el primer evento
func (client *AMIClient) run()  {

	go func(){

		for {

			data, _ := client.conn.ReadMIMEHeader()

			if ev, err := newEvent(&data); err == nil {
				client.Events <- *ev

			}
			if response, err := newResponse(&data); err == nil {
				fmt.Printf("\tDATA:%v\n", *response)
				client.response <- *response
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

func newEvent(data *textproto.MIMEHeader) (*AMIEvent, error)  {
	if data.Get("Event") == "" {
		return nil, errors.New("Not Event")
	}
	ev := &AMIEvent{data.Get("Event"), strings.Split(data.Get("Privilege"), ","), make(map[string]string)}
	for k,v := range *data {
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
	if strings.Contains(label, "Asterisk Call Manager")  != true {
		return nil, errNotAMI
	}

	client := &AMIClient{conn, make(chan AMIResponse), make(chan AMIEvent, 100)}
	client.run()
	return client, nil
}

