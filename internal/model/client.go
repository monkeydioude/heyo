package model

import (
	"time"

	"github.com/monkeydioude/heyo/pkg/rpc"
	"google.golang.org/grpc"
)

// Client model is the representation of a client, in the server's perspective
type Client struct {
	Event          string
	Uuid           string
	Name           string
	SubscribedAt   time.Time
	ResponseSender func(res *rpc.Message) error
	MessageChan    chan *rpc.Message
}

func (cl *Client) Send(in *rpc.Message) error {
	if in == nil {
		return ErrNilParameter
	}
	cl.MessageChan <- in
	return nil
}

type ClientFactory struct {
	timeFn func() time.Time
}

func NewFactory() ClientFactory {
	return ClientFactory{
		timeFn: func() time.Time { return time.Now() },
	}
}

func (cf *ClientFactory) NewFromSubscription(
	sub *rpc.Subscriber,
	res grpc.ServerStreamingServer[rpc.Message],
) Client {
	return Client{
		Event:          sub.Event,
		ResponseSender: res.Send,
		Uuid:           sub.ClientId,
		SubscribedAt:   cf.timeFn(),
		MessageChan:    make(chan *rpc.Message, 100),
		Name:           sub.Name,
	}
}
