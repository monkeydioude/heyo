package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/monkeydioude/heyo/internal/consts"
	"github.com/monkeydioude/heyo/internal/handler/server"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func grpcRouting(ctx context.Context, s *grpc.Server) {
	rpc.RegisterBrokerServer(s, server.NewHeyoServer(ctx))
}

func grpcServer(ctx context.Context, runGroup *run.Group) {
	port := consts.RPCPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	ctx, cancelFn := context.WithCancel(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("[ERR ] failed to listen: %v", err)
	}
	// creds, err := credentials.NewServerTLSFromFile("./certs/localhost.crt", "./certs/localhost.key")
	// if err != nil {
	// 	log.Fatalf("[ERR ] Invalid creds: %v", err)
	// }
	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
	)
	grpcRouting(ctx, server)
	runGroup.Add(func() error {
		log.Println("[INFO] RPC starting on", lis.Addr())
		return server.Serve(lis)
	}, func(err error) {
		log.Println("[INFO] stopping RPC server")
		cancelFn()
		server.GracefulStop()
		server.Stop()
	})
}
