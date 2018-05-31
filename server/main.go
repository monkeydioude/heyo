package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	sc "github.com/monkeydioude/schampionne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":9393"

type server struct {
	listeners map[string][]chan *sc.Rumor
}

// Ack returns a pointer to a sc.Ack struct
func Ack(m string, code sc.AckCode) *sc.Ack {
	return &sc.Ack{
		Message: m,
		Code:    code,
	}
}

func (s *server) Whisper(c context.Context, r *sc.Rumor) (*sc.Ack, error) {
	if _, ok := s.listeners[r.Type]; !ok {
		return Ack("No listener", sc.AckCode_NO_LISTENER), fmt.Errorf("No listeners for %s", r.Type)
	}

	fmt.Printf("Receiving Rumor %+v\n", r)

	for _, c := range s.listeners[r.Type] {
		c <- r
	}

	return Ack("Ok", sc.AckCode_OK), nil
}

func (s *server) Listen(l *sc.Listener, stream sc.Broker_ListenServer) error {
	cr := make(chan *sc.Rumor)
	s.listeners[l.Type] = append(s.listeners[l.Type], cr)

	for {
		r := <-cr
		if r == nil {
			log.Printf("[INFO] Listener for %s ended\n", l.Type)
			break
		}
		stream.Send(r)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sc.RegisterBrokerServer(s, &server{
		listeners: make(map[string][]chan *sc.Rumor),
	})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
