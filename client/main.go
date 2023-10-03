// Package main implements a greeting client service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/peer"

	_ "google.golang.org/grpc/xds/experimental"
)

var (
	serverAddr = flag.String("serverAddr", "localhost:50051", "Address of the server with port")
	greetName  = flag.String("greetName", "world", "Name to send greeting")

	connTimeout     = flag.Duration("connTimeout", 10*time.Second, "Duration for establishing grpc.ClientConn, format: 100ms, or 10s, or 2h")
	requestTimeout  = flag.Duration("requestTimeout", time.Second, "Duration for each RPC, format: 100ms, or 10s, or 2h")
	runningDuration = flag.Duration("duration", time.Hour, "Total duration the program runs, format: 100ms, or 10s, or 2h")

	useXDS = flag.Bool("useXDS", true, "Enable or disable xds")
)

func main() {
	flag.Parse()

	ctx, terminate := context.WithTimeout(context.Background(), *connTimeout)
	defer terminate()

	targetAddr := *serverAddr
	if *useXDS {
		targetAddr = "xds-experimental:///" + targetAddr
	}

	connection, err := grpc.DialContext(ctx, targetAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to establish connection: %v", err)
	}
	defer connection.Close()

	client := pb.NewGreeterClient(connection)

	termination := time.After(*runningDuration)

	for {
		requestCtx, terminateReq := context.WithTimeout(context.Background(), *requestTimeout)
		serverPeer := new(peer.Peer)
		response, err := client.SayHello(requestCtx, &pb.HelloRequest{Name: *greetName}, grpc.Peer(serverPeer))
		if err != nil {
			terminateReq()
			log.Fatalf("Greeting failed: %v", err)
		}
		terminateReq()
		fmt.Printf("Received Greeting: %s, from %v\n", response.GetMessage(), serverPeer.Addr)

		select {
		case <-termination:
			return
		default:
		}
		time.Sleep(time.Second)
	}
}
