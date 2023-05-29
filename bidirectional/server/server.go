package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

var clientConn net.Conn
var connected = false

type Arith int

func (t *Arith) Mul(ctx context.Context, args *example.Args, reply *example.Reply) error {

	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	// time.Sleep(10 * time.Second)
	reply.C = args.A * args.B
	connected = true
	return nil
}

func (t *Arith) Mul2(ctx context.Context, args *example.Args, reply *example.Reply) error {

	clientConn = ctx.Value(server.RemoteConnContextKey).(net.Conn)
	reply.C = args.A * args.B
	connected = true
	return nil
}

func main() {
	flag.Parse()

	ln, _ := net.Listen("tcp", ":9981")
	go http.Serve(ln, nil)

	s := server.NewServer()
	//s.RegisterName("Arith", new(example.Arith), "")
	s.Register(new(Arith), "")
	go s.Serve("tcp", *addr)

	for !connected {
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("start to send messages to %s\n", clientConn.RemoteAddr().String())
	for i := 1; i < 100000; i++ {
		if clientConn != nil {
			go func(i int) {
				// err := s.SendMessage(clientConn, "test_service_path", "test_service_method", nil, []byte("abcde"))
				err := s.SendMessage(clientConn, "test_service_path", "test_service_method", nil, []byte(strconv.Itoa(i)))
				if err != nil {
					fmt.Printf("failed to send messsage to %s: %v\n", clientConn.RemoteAddr().String(), err)
					clientConn = nil
				}
			}(i)

		}
		// time.Sleep(time.Second)
	}
	select {}
}
