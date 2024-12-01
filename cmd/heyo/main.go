package main

import (
	"context"
	"log"
	"syscall"

	"github.com/oklog/run"
)

func main() {
	runGroup := run.Group{}
	ctx := context.TODO()
	restServer(&runGroup)
	grpcServer(ctx, &runGroup)
	// Signals handling, for graceful stop
	runGroup.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	if err := runGroup.Run(); err != nil {
		log.Fatal(err)
	}
}
