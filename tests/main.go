package main

import (
	"context"
	"fmt"
	"io"
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

	client := sc.NewBrokerClient(conn)
	listener := &sc.Listener{Type: "test"}

	stream, err := client.Listen(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	err = stream.Send(listener)
	if err != nil {
		log.Fatalln(err)
	}

	defer stream.CloseSend()
	defer conn.Close()

	for {
		rumor, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("Wesh alors :(")
			break
		}

		if err != nil {
			log.Fatalf("Heee?: %s \n", err)
		}

		fmt.Printf("HELLOOOO %+v\n", rumor)
	}
}
