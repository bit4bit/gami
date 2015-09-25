package gami

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"sync"
	"testing"
	"time"
)

type amiMockAction func(params textproto.MIMEHeader) map[string]string

//amiServer for mocking Asterisk AMI
type amiServer struct {
	Addr          string
	actionsMocked map[string]amiMockAction
	listener      net.Listener
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

	//example mocking login of asterisk
	srv.Mock("Login", func(params textproto.MIMEHeader) map[string]string {
		return map[string]string{
			"Response": "OK",
			"ActionID": params.Get("Actionid"),
		}
	})
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
	workers := 5

	wg := &sync.WaitGroup{}
	for ti := tests; ti > 0; ti-- {
		resWorkers := make(chan (<-chan *AMIResponse), workers)
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				chres, err := ami.AsyncAction("Test", nil)
				if err != nil {
					t.Error(err)
				}
				resWorkers <- chres
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(resWorkers)
		}()

		for resp := range resWorkers {
			select {
			case <-time.After(time.Second * 5):
				t.Fatal("asyncAction locked")
			case <-resp:
			}
		}

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
	srv := &amiServer{Addr: listener.Addr().String(),
		listener:      listener,
		actionsMocked: make(map[string]amiMockAction)}
	go srv.do(listener)
	return srv
}

//Mock
func (c *amiServer) Mock(action string, cb amiMockAction) {
	c.actionsMocked[action] = cb
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

		go func(conn *textproto.Conn) {
			defer conn.Close()

			for {
				header, err := conn.ReadMIMEHeader()
				if err != nil {
					return
				}
				var output bytes.Buffer

				time.AfterFunc(time.Millisecond*time.Duration(rand.Intn(1000)), func() {

					if _, ok := c.actionsMocked[header.Get("Action")]; ok {
						rvals := c.actionsMocked[header.Get("Action")](header)
						for k, vals := range rvals {
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
