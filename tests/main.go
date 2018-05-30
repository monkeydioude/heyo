package main

import (
	"fmt"
	"io"
	"log"

	sc "github.com/monkeydioude/schampionne"
	"golang.org/x/net/context"
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
	listener := &sc.Listener{Type: "test"}
	stream, err := client.Listen(context.Background(), listener)

	if err != nil {
		log.Fatalln(err)
	}

	for {
		rumor, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("Wesh alors :(")
			break
		}

		if err != nil {
			log.Fatalf("%s zbrah\n", err)
		}

		fmt.Printf("%+v\n", rumor)
	}
	stream.CloseSend()
}
