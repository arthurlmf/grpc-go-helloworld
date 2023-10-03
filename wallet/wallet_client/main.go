// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	walletpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet"

	_ "google.golang.org/grpc/xds/experimental"
)

var (
	address = flag.String("address", "localhost:50051", "server address, with port")
	name    = flag.String("name", "world", "name to greet")

	dialTimeout = flag.Duration("dialTimeout", 10*time.Second, "timeout for creating grpc.ClientConn, format: 100ms, or 10s, or 2h")
	rpcTimeout  = flag.Duration("rpcTimeout", time.Second, "timeout for each RPC, format: 100ms, or 10s, or 2h")

	watch = flag.Bool("watch", false, "streaming watch or not")
	totalTime = flag.Duration("time", time.Hour, "total time the binary runs, format: 100ms, or 10s, or 2h")
	gap = flag.Duration("gap", time.Second, "sleep between RPCs")
)

func fetchB(end <-chan time.Time, c walletpb.WalletClient) {
	for {
		var header metadata.MD
		ctx, cancel := context.WithTimeout(context.Background(), *rpcTimeout)
		defer cancel()
		r, err := c.FetchBalance(ctx, &walletpb.BalanceRequest{}, grpc.Header(&header))
		if err != nil {
			log.Fatalf("could not fetch: %v", err)
		}
		fmt.Printf("Fetch: %v, from %v\n", r, header["hostname"])

		select {
		case <-end:
			return
		default:
		}
		time.Sleep(*gap)
	}
}

func watchB(end <-chan time.Time, c walletpb.WalletClient) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), *rpcTimeout)
		defer cancel()
		s, err := c.WatchBalance(ctx, &walletpb.BalanceRequest{})
		if err != nil {
			log.Fatalf("could not fetch: %v", err)
		}
		header, _ := s.Header()
		r, _ := s.Recv()
		fmt.Printf("Watch: %v, from %v\n", r, header["hostname"])


		select {
		case <-end:
			return
		default:
		}
		time.Sleep(*gap)
	}
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), *dialTimeout)
	defer cancel()

	addr := *address

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := walletpb.NewWalletClient(conn)

	end := time.After(*totalTime)

	if !*watch {
		fetchB(end, c)
		return
	}
	watchB(end, c)
}
