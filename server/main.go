package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// Import the generated protobuf package
	pb "xp-go-grpc/proto"

	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received request for name=%s", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
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
	pb.RegisterGreeterServer(s, &server{})

	fmt.Println("Server listening at port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
