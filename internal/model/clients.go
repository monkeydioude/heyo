package model

import (
	"errors"
	"fmt"
	"slices"

	"github.com/monkeydioude/heyo/pkg/rpc"
)

type Clients map[string][]*Client

var ErrCouldNotAddClient = errors.New("could not add client")
var ErrCouldNotRemoveClient = errors.New("could not add client")
var ErrNilParameter = errors.New("nil parameter")
var ErrTryReachingUnknownEvent = errors.New("try to reach an unknown event")

func NewClients() Clients {
	return make(map[string][]*Client, 0)
}

func (cls Clients) Add(client *Client) error {
	if client == nil {
		return errors.Join(ErrCouldNotAddClient, ErrNilParameter)
	}
	slice, ok := cls[client.Event]
	if !ok {
		slice = make([]*Client, 0)
	}
	slice = append(slice, client)
	cls[client.Event] = slice
	return nil
}

func (cls Clients) Remove(client *Client) error {
	if client == nil {
		return errors.Join(ErrCouldNotRemoveClient, ErrNilParameter)
	}
	clientSlice, ok := cls[client.Event]
	if !ok {
		return nil
	}
	clientSlice = slices.DeleteFunc(clientSlice, func(cl *Client) bool {
		return cl != nil && cl.Event == client.Event && cl.Uuid == client.Uuid
	})
	cls[client.Event] = clientSlice
	return nil
}

func (cls Clients) TotalLen() int {
	n := 0
	for _, subs := range cls {
		n += len(subs)
	}
	return n
}

func (cls Clients) Len(event string) int {
	subs, ok := cls[event]
	if !ok {
		return 0
	}
	return len(subs)
}

func (cls Clients) Send(in *rpc.Message) error {
	if in == nil {
		return ErrNilParameter
	}
	clients, ok := cls[in.Event]
	if !ok {
		return fmt.Errorf("%s: %w", in.Event, ErrTryReachingUnknownEvent)
	}
	for _, client := range clients {
		client.Send(in)
	}
	return nil
}
