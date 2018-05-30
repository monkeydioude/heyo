package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	sc "github.com/monkeydioude/schampionne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":9393"

type server struct{}

func (s *server) Whisper(context.Context, *sc.Rumor) (*sc.Ack, error) {
	return nil, nil
}

func (s *server) Listen(l *sc.Listener, stream sc.Broker_ListenServer) error {
	// for {
	stream.Send(&sc.Rumor{
		Type:    "test",
		Message: "pouet",
	})
	// }
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sc.RegisterBrokerServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
