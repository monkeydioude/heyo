package boot

import (
	"context"
	"fmt"
	"log"

	"github.com/monkeydioude/heyo/internal/consts"
	"github.com/monkeydioude/heyo/internal/service/client"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/monkeydioude/heyo/pkg/tiger/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getRPCClient() rpc.BrokerClient {
	cl, err := grpc.NewClient(
		fmt.Sprintf("[::]:%s", consts.RPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(err)
	assert.NotNil(cl)
	log.Printf("[INFO] connecting to port [::]:%s\n", consts.RPCPort)
	return rpc.NewBrokerClient(cl)
}

func BootClient(ctx context.Context) client.Client {
	return client.New(ctx, getRPCClient())
}
