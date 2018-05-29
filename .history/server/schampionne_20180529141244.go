package main

import (
	"context"
	"log"
	"net"

	"github.com/monkeydioude/schampionne/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Whisper(context.Context, *grpc.Rumor) (*grpc.Response, error) {
	return nil, nil
}

func (s *server) Listen(context.Context, *grpc.Listener) (*grpc.Rumor, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", "9393")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
