package client

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/monkeydioude/heyo/pkg/rpc"
	"google.golang.org/grpc"
)

func StreamFetchMessage(
	stream grpc.ServerStreamingClient[rpc.Message],
) (*rpc.Message, error) {
	msg := rpc.Message{}
	err := stream.RecvMsg(&msg)
	if err != nil {
		return &msg, fmt.Errorf("%w: %w", ErrMessageReception, err)
	}
	return &msg, nil
}

func Listener(
	stream grpc.ServerStreamingClient[rpc.Message],
	onMessageReceived func(*rpc.Message) error,
	mutex *sync.Mutex,
) error {
	defer stream.CloseSend()
	for {
		msg := rpc.Message{}
		err := stream.RecvMsg(&msg)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrMessageReception, err)
		}
		mutex.Lock()
		if err := onMessageReceived(&msg); err != nil {
			if errors.Is(err, ErrListenFatalErr) {
				return fmt.Errorf("%w: %w", ErrListenFatalErr, err)
			}
			log.Printf("[WARN] %s\n", err.Error())
		}
		mutex.Unlock()
	}
}
