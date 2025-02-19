package main

import (
	"context"
	"errors"
	"flag"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	s.RegisterName("Arith", new(example.Arith), "")
	s.AuthFunc = auth
	s.Serve("reuseport", *addr)
}

func auth(ctx context.Context, req *protocol.Message, token string) error {

	if token == "bearer abcdefg1234567" {
		return nil
	}

	return errors.New("invalid token")
}
