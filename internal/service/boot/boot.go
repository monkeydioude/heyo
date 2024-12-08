package boot

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/monkeydioude/heyo/internal/consts"
	"github.com/monkeydioude/heyo/internal/service/client"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/monkeydioude/heyo/pkg/tiger/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func getRPCClient() rpc.BrokerClient {
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // Skip verification for testing; remove this in production
	})

	port := consts.RPCPort
	if os.Getenv("SERVER_PORT") != "" {
		port = os.Getenv("SERVER_PORT")
	}
	cl, err := grpc.NewClient(
		fmt.Sprintf("[::]:%s", port),
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Second, // Ping every 10 seconds
			Timeout:             5 * time.Second, // Wait 5 seconds for a ping response
			PermitWithoutStream: true,            // Allow pings even without active streams
		}),
	)
	assert.NoError(err)
	assert.NotNil(cl)
	log.Printf("[INFO] connecting to port [::]:%s\n", port)
	return rpc.NewBrokerClient(cl)
}

func BootClient(ctx context.Context) client.Client {
	return client.New(ctx, getRPCClient())
}
