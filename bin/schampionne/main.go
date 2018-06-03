package main

import (
	"errors"
	"log"
	"net"

	"golang.org/x/net/context"

	sc "github.com/monkeydioude/schampionne"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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
	log.Printf("[INFO] Receiving Rumor %+v\n", r)

	if _, ok := s.listeners[r.Type]; !ok || len(s.listeners) == 0 {
		log.Printf("[INFO] No listener for Type \"%s\"\n", r.Type)
		return Ack("No listener", sc.AckCode_NO_LISTENER), errors.New("No listener")
	}

	for _, c := range s.listeners[r.Type] {
		c <- r
	}

	return Ack("Ok", sc.AckCode_OK), nil
}

func (s *server) Listen(stream sc.Broker_ListenServer) error {
	l, err := stream.Recv()
	if err != nil {
		log.Printf("[WARN] %s\n", err)
	}

	p, ok := peer.FromContext(stream.Context())
	if ok == false {
		log.Println("[WARN] Could not fetch peer Context")
	}

	log.Printf("[INFO] Incoming Listener for Type \"%s\" %+v\n", l.Type, p)

	cr := make(chan *sc.Rumor)
	s.listeners[l.Type] = append(s.listeners[l.Type], cr)

	go func() {
		<-stream.Context().Done()
		cr <- nil
	}()

	for {
		r := <-cr
		if r == nil {
			log.Printf("[INFO] Listener for Type \"%s\" disconnected %+v\n", l.Type, p)
			delete(s.listeners, l.Type)
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

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
