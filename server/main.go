// Package main provides a server for the Greeting service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	serverPort     = flag.Int("serverPort", 50051, "The main server port")
	healthChkPort  = flag.Int("healthChkPort", 49940, "Port for health checks")
	currentHost, _ = os.Hostname()
)

// greetingServer represents the implementation of helloworld.GreeterServer.
type greetingServer struct{}

// Greet implements helloworld.GreeterServer.
func (g *greetingServer) Greet(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("Received from: %v\n", request.GetName())
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, this is %s", request.GetName(), currentHost)}, nil
}

func main() {
	flag.Parse()

	go func() {
		http.HandleFunc("/", func(response http.ResponseWriter, req *http.Request) {
			fmt.Println("Health check initiated")
			fmt.Fprintf(response, "Hello")
		})
		http.ListenAndServe(fmt.Sprintf(":%d", *healthChkPort), nil)
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("Error listening on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &greetingServer{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
