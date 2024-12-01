package main

import (
	"context"
	"errors"
	"log"
)

func main() {
	ctx := context.TODO()
	explorer := NewExplorer(ctx)
	explorer.setupReceiverEvents()
	for {
		in, err := buildMessage(explorer.GetCtx(), explorer.client.Uuid)
		if err != nil {
			if errors.Is(err, ErrUpstreamClosed) {
				log.Fatalf("[ERR ] stream closed: %s\n", err)
			}
			log.Printf("[WARN] %s\n", err.Error())
		}
		if err := explorer.sendMessage(&in); err != nil {
			log.Printf("[ERR ] %s\n", err)
		}
	}
}
