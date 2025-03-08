package main

import (
	"fmt"
	"log"
	"net"

	// Import the generated protobuf package
	pb "xp-go-grpc-server/proto"

	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedAccountsServiceServer
}

func main() {
	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Create a gRPC server object
	s := grpc.NewServer()

	// Register our service implementation with the gRPC server
	pb.RegisterAccountsServiceServer(s, &server{})

	fmt.Println("Server listening at port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
