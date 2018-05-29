package main

import (
	"context"
	"log"
	"net"

	sc "github.com/monkeydioude/schampionne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Whisper(context.Context, *sc.Rumor) (*sc.Response, error) {
	return nil, nil
}

func (s *server) Listen(context.Context, *sc.Listener) (*sc.Rumor, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", "9393")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRumorServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
