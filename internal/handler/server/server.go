package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/monkeydioude/heyo/internal/model"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"google.golang.org/grpc"
)

type HeyoServer struct {
	rpc.UnimplementedBrokerServer
	clientFactory *model.ClientFactory
	clients       model.Clients
	mutex         sync.Mutex
	ctx           context.Context
}

func (hs *HeyoServer) Enqueue(ctx context.Context, in *rpc.Message) (*rpc.Ack, error) {
	if in == nil {
		return nil, fmt.Errorf("Enqueue: %w", ErrNilIncomingMessage)
	}
	log.Printf("[INFO] '%s': from client '%s:%s', message_id '%s', payload: %s\n", in.Event, in.ClientName, in.ClientId, in.MessageId, in.Data)
	err := hs.clients.Send(in)
	if err != nil {
		code := rpc.AckCode_INTERNAL_ERROR
		if errors.Is(err, model.ErrTryReachingUnknownEvent) {
			code = rpc.AckCode_UNKNOWN_EVENT
		}
		return &rpc.Ack{
			Data: "nok",
			Code: code,
		}, err
	}
	return &rpc.Ack{
		Data: "ok",
		Code: rpc.AckCode_OK,
	}, nil
}

func (hs *HeyoServer) Subscription(
	subscriber *rpc.Subscriber,
	res grpc.ServerStreamingServer[rpc.Message],
) error {
	hs.mutex.Lock()
	client := hs.clientFactory.NewFromSubscription(subscriber, res)
	hs.clients.Add(&client)
	log.Printf("[INFO] '%s:%s' subscribed to '%s'\n", client.Name, client.Uuid, client.Event)
	log.Printf("[INFO] '%s' now has %d subscribers", subscriber.Event, hs.clients.Len(subscriber.Event))
	hs.mutex.Unlock()

	// Clean up on client disconnect
	defer hs.ClientDisconnect(&client, subscriber)

	for {
		select {
		case msg := <-client.MessageChan:
			log.Printf("[INFO] '%s:%s': sending message_id '%s' to client '%s'\n", client.Name, client.Event, msg.MessageId, client.Uuid)
			// Send the message to the client
			if err := res.Send(msg); err != nil {
				log.Printf("[ERR ] error sending message to '%s' of '%s': %v\n", client.Uuid, client.Event, err)
				return err
			}
		case <-hs.ctx.Done():
			return nil
		case <-res.Context().Done():
			// Client has disconnected
			return nil
		}
	}
}

func (hs *HeyoServer) ClientDisconnect(client *model.Client, subscriber *rpc.Subscriber) {
	hs.mutex.Lock()
	hs.clients.Remove(client)
	hs.mutex.Unlock()
	if client.MessageChan != nil {
		close(client.MessageChan)
	}
	log.Printf("[INFO] '%s:%s' disconnected from '%s'\n", client.Name, client.Uuid, client.Event)
	log.Printf("[INFO] '%s' now has %d subscribers", subscriber.Event, hs.clients.Len(subscriber.Event))
}

func NewHeyoServer(ctx context.Context) *HeyoServer {
	cf := model.NewFactory()
	return &HeyoServer{
		clientFactory: &cf,
		clients:       model.NewClients(),
		ctx:           ctx,
	}
}
