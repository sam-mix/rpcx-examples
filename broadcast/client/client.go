package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	addr1 = flag.String("addr1", "tcp@localhost:8972", "server1 address")
	addr2 = flag.String("addr2", "tcp@localhost:9981", "server2 address")
)

func main() {
	flag.Parse()
	for {
		process()
	}
}

func process() {
	d, _ := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	// for {
	reply := &example.Reply{}
	err := xclient.Broadcast(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	if reply.C == 200 {
		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	}
	// break
	// time.Sleep(1e9)
	// }
}
