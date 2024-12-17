package boot

import (
	"context"
	"log"
	"os"

	"github.com/monkeydioude/heyo/internal/consts"
	"github.com/monkeydioude/heyo/internal/service/client"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/monkeydioude/heyo/pkg/tiger/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getRPCClient() rpc.BrokerClient {
	// creds := credentials.NewTLS(&tls.Config{
	// 	InsecureSkipVerify: true, // Skip verification for testing; remove this in production
	// })

	addr := consts.RPCAddr
	if os.Getenv("SERVER_ADDR") != "" {
		addr = os.Getenv("SERVER_ADDR")
	}
	cl, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// 	Time:                10 * time.Second, // Ping every 10 seconds
		// 	Timeout:             5 * time.Second,  // Wait 5 seconds for a ping response
		// 	PermitWithoutStream: true,             // Allow pings even without active streams
		// }),
	)
	assert.NoError(err)
	assert.NotNil(cl)
	log.Printf("[INFO] connecting to %s\n", addr)
	return rpc.NewBrokerClient(cl)
}

func BootClient(ctx context.Context) client.Client {
	return client.New(ctx, getRPCClient())
}
