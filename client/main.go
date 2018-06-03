package main

import (
	"fmt"
	"log"

	sc "github.com/monkeydioude/schampionne"
)

const (
	address = "localhost:9393"
)

func main() {
	client := sc.NewClient(address)

	defer client.Close()

	ack, err := client.Whisper("test", "pouet")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ack received: m:(%s), c:(%d)\n", ack.GetMessage(), ack.GetCode())
}
