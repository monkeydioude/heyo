package main

import (
	"context"
	"fmt"
	"log"

	sc "github.com/monkeydioude/schampionne"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9393"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer conn.Close()

	client := sc.NewBrokerClient(conn)

	r := &sc.Rumor{Type: "test", Message: "test"}

	fmt.Printf("SENDING: %+v, ", r)

	ack, _ := client.Whisper(context.Background(), r)

	fmt.Printf("Ack received: m:(%s), c:(%d)\n", ack.Message, ack.Code)
}
