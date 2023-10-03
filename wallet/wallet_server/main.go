// server ...
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	walletpb "google.golang.org/grpc/grpc-wallet/grpc/examples/wallet"
	"google.golang.org/grpc/metadata"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	serverID    = flag.String("id", "not-set", "")
	hostname, _ = os.Hostname()
)

type server struct{}

func (s *server) FetchBalance(ctx context.Context, _ *walletpb.BalanceRequest) (*walletpb.BalanceResponse, error) {
	grpc.SetHeader(ctx, metadata.Pairs("hostname", hostname, "hostname", *serverID))
	return &walletpb.BalanceResponse{Balance: 123}, nil
}

func (s *server) WatchBalance(_ *walletpb.BalanceRequest, stream walletpb.Wallet_WatchBalanceServer) error {
	stream.SetHeader(metadata.Pairs("hostname", hostname, "hostname", *serverID))
	stream.Send(&walletpb.BalanceResponse{Balance: 456})
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	walletpb.RegisterWalletServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
