package gami

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"reflect"
	"testing"
	"time"
)

//amiServer for mocking Asterisk AMI
type amiServer struct {
	Addr     string
	listener net.Listener
}

func TestLogin(t *testing.T) {
	srv := newAmiServer()
	defer srv.Close()
	ami, err := Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	go ami.Run()
	defer ami.Close()
	defaultInstaller(t, ami)
	ami.Login("admin", "admin")
}

func TestMultiAsyncActions(t *testing.T) {
	srv := newAmiServer()
	defer srv.Close()
	ami, err := Dial(srv.Addr)
	if err != nil {
		t.Fatal(err)
	}
	go ami.Run()
	defer ami.Close()
	defaultInstaller(t, ami)

	tests := 10
	workers := 3

	for ti := tests; ti > 0; ti-- {
		resWorkers := make(chan (<-chan *AMIResponse), workers)
		for i := 0; i < workers; i++ {
			chres, err := ami.AsyncAction("Test", nil)
			if err != nil {
				t.Error(err)
			}
			resWorkers <- chres
		}
		close(resWorkers)
		done := make(chan bool)
		go func() {
			select {
			case <-time.After(time.Second * 2):
				t.Error("asyncAction locked")
			case <-done:
				done <- true
				break
			}
		}()
		for res := range resWorkers {
			<-res
		}
		done <- true
		<-done
	}

}

func defaultInstaller(t *testing.T, ami *AMIClient) {
	go func() {
		for {
			select {
			//handle network errors
			case err := <-ami.NetError:
				t.Error("Network Error:", err)
			case err := <-ami.Error:
				t.Error("error:", err)
			//wait events and process
			case <-ami.Events:
				//t.Log("Event:", *ev)
			}
		}
	}()
}

func newAmiServer() *amiServer {
	addr := "localhost:0"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv := &amiServer{Addr: listener.Addr().String(), listener: listener}
	go srv.do(listener)
	return srv
}

//MockLogin it's a example for mocking responses
func (c *amiServer) MockLogin(params textproto.MIMEHeader) map[string]string {
	println("Llamado login")
	return map[string]string{
		"Response": "OK",
		"ActionID": params.Get("Actionid"),
	}
}

func (c *amiServer) do(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		fmt.Fprintf(conn, "Asterisk Call Manager\r\n")
		tconn := textproto.NewConn(conn)
		//install event HeartBeat
		go func(conn *textproto.Conn) {
			for now := range time.Tick(time.Second) {
				fmt.Fprintf(conn.W, "Event: HeartBeat\r\nTime: %d\r\n\r\n",
					now.Unix())
			}
		}(tconn)

		rval := reflect.ValueOf(c)
		go func(conn *textproto.Conn) {
			defer conn.Close()

			for {
				header, err := conn.ReadMIMEHeader()
				if err != nil {
					return
				}
				var output bytes.Buffer

				time.AfterFunc(time.Millisecond*time.Duration(rand.Intn(1000)), func() {
					fnc := rval.MethodByName("Mock" + header.Get("Action"))
					if fnc.IsValid() {
						rvals := fnc.Call([]reflect.Value{reflect.ValueOf(header)})
						ival := rvals[0].Interface()
						for k, vals := range ival.(map[string]string) {
							fmt.Fprintf(&output, "%s: %s\r\n", k, vals)
						}
						output.WriteString("\r\n")

						err := conn.PrintfLine(output.String())
						if err != nil {
							panic(err)
						}
					} else {
						//default response
						fmt.Fprintf(&output, "Response: TEST\r\nActionID: %s\r\n\r\n",
							header.Get("Actionid"))
						err := conn.PrintfLine(output.String())
						if err != nil {
							panic(err)
						}
					}
				})

			}
		}(tconn)
	}
}

func (c *amiServer) Close() {
	c.listener.Close()
}
