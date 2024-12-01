package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/monkeydioude/heyo/pkg/datatype/mapvec"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"google.golang.org/grpc"
)

type Streams = mapvec.MapVec[string, grpc.ServerStreamingClient[rpc.Message]]

type Client struct {
	RpcClient   rpc.BrokerClient
	streams     Streams
	Uuid        string
	ctx         context.Context
	ctxCancelFn context.CancelFunc
	mutex       sync.Mutex
}

func New(ctx context.Context, RpcClient rpc.BrokerClient) Client {
	_ctx, _cancelFn := context.WithCancel(ctx)
	return Client{
		RpcClient:   RpcClient,
		Uuid:        uuid.NewString(),
		streams:     make(Streams),
		ctx:         _ctx,
		ctxCancelFn: _cancelFn,
	}
}

// func (cl *Client) Handshake(
// 	stream grpc.ServerStreamingClient[rpc.Message],
// ) error {
// 	err := async.Timeout(5*time.Second, func() error {
// 		log.Printf("[INFO] starting handshake")
// 		msg, err := StreamFetchMessage(stream)
// 		if err != nil {
// 			return fmt.Errorf("%w: %w", ErrHandshake, err)
// 		}
// 		hs := client.Handshake{}
// 		err = json.Unmarshal([]byte(msg.Data), &hs)
// 		if err != nil {
// 			return err
// 		}
// 		cl.mutex.Lock()
// 		cl.Uuid = hs.ClientUuid
// 		log.Printf("[INFO] handshake SUCCESS, setting client uuid to %s", cl.Uuid)
// 		cl.mutex.Unlock()
// 		return nil
// 	}, func() error {
// 		cl.mutex.Lock()
// 		cl.Uuid = uuid.NewString()
// 		log.Printf("[WARN] handshake FAIL, self-setting client uuid to %s", cl.Uuid)
// 		cl.mutex.Unlock()
// 		return nil
// 	})

// 	if err != nil {
// 		return fmt.Errorf("%w: %w", ErrHandshake, err)
// 	}
// 	return nil
// }

func (cl *Client) MakeSubscription(event string) (grpc.ServerStreamingClient[rpc.Message], error) {
	sub := rpc.Subscriber{
		Event:      event,
		ClientUuid: cl.Uuid,
	}
	stream, err := cl.RpcClient.Subscription(cl.ctx, &sub)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (cl *Client) Listen(
	event string,
	onMessageReceived func(*rpc.Message) error,
) error {
	stream, err := cl.MakeSubscription(event)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSubscriptionFailed, err)
	}
	log.Printf("[INFO] subscribed to '%s'\n", event)
	cl.streams.Add(event, stream)
	// err = cl.Handshake(stream)
	// if err != nil {
	// 	log.Printf("[ERR ] event '%s': %s \n", event, err.Error())
	// 	cl.Close()
	// }
	go func() {
		if err := Listener(stream, onMessageReceived, &cl.mutex); err != nil {
			log.Printf("[ERR ] event '%s': %s \n", event, err.Error())
			cl.Close()
		}
	}()
	return nil
}

func (cl *Client) Send(event string, msg *rpc.Message) error {
	msg.ClientUuid = cl.Uuid
	ctx, cancelFn := context.WithTimeout(cl.ctx, 5*time.Second)
	defer cancelFn()
	ack, err := cl.RpcClient.Enqueue(ctx, msg)
	if err != nil {
		return err
	}
	if ack == nil {
		return fmt.Errorf("%w: %w", ErrListenFatalErr, ErrNilPointer)
	}
	switch ack.Code {
	case rpc.AckCode_INTERNAL_ERROR:
		log.Printf("[ERR ] event '%s' experience an internal error: %s\n", event, ack.Data)
	case rpc.AckCode_NO_LISTENER:
		log.Printf("[WARN] event '%s' has no listeners\n", event)
	case rpc.AckCode_QUEUE_FULL:
		log.Printf("[ERR ] event '%s' queue is full\n", event)
	case rpc.AckCode_UNKNOWN_EVENT:
		log.Printf("[WARN] event '%s' is unknown\n", event)
	case rpc.AckCode_OK:
		log.Printf("[INFO] event '%s' OK: %s\n", event, ack.Data)
	}
	return nil
}

func (cl *Client) GetCtx() context.Context {
	return cl.ctx
}

func (cl *Client) Close() error {
	cl.ctxCancelFn()
	return nil
}
