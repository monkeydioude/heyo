package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/monkeydioude/heyo/internal/service/boot"
	"github.com/monkeydioude/heyo/internal/service/client"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/monkeydioude/heyo/pkg/tiger/assert"
)

type Mode string

const (
	ModeEnqueue      Mode = "enq"
	ModeSubscription Mode = "sub"
)

func ModeMatches(m string) bool {
	return m == string(ModeEnqueue) || m == string(ModeSubscription)
}

func Into(m string) Mode {
	if !ModeMatches(m) {
		log.Fatalf("mode '%s' unknown", m)
	}
	return Mode(m)
}

func (m Mode) IsEnq() bool {
	return m == ModeEnqueue
}

func (m Mode) Do(
	ctx context.Context,
	client *client.Client,
	event string,
) error {
	assert.NotNil(client)
	assert.NotEmpty(event)
	err := errors.New("unexpected mode")
	switch m {
	case ModeEnqueue:
		for i := 0; i < flag.NArg(); i++ {
			msg := flag.Arg(i)
			assert.NotEmpty(msg)
			log.Println("sending message", msg)
			_, err = client.RpcClient.Enqueue(ctx, &rpc.Message{
				Event:     event,
				Data:      msg,
				MessageId: uuid.NewString(),
				ClientId:  client.Uuid,
			})
		}
	case ModeSubscription:
		ctx, cancel := context.WithCancel(client.GetCtx())
		msgChan := make(chan string, 100)
		defer cancel()
		go client.Listen(event, func(m *rpc.Message) error {
			if m == nil {
				return errors.New("nil message")
			}
			msgChan <- m.Data
			return nil
		})
		for {
			select {
			case <-ctx.Done():
				return nil
			case msg := <-msgChan:
				fmt.Println(msg)
			}
		}
	}
	return err
}

type Event = string
type Message = string

func getParams() (Mode, Event) {
	mode := flag.String("mode", "enq", "'enq' Enqueue\n 'sub'")
	event := flag.String("event", "", "the event to enqueue or subscribe to")
	flag.Parse()
	assert.NotNilNorEmpty(mode, errors.New("must specify a mode"))
	assert.NotNilNorEmpty(event, errors.New("must specify an event"))
	m := Into(*mode)
	if m.IsEnq() && flag.NArg() == 0 {
		log.Fatal("need a message to enqueue")
	}
	return m, *event
}

func main() {
	mode, event := getParams()
	assert.NotEmpty(mode)
	assert.NotEmpty(event)

	ctx := context.TODO()
	client := boot.BootClient(ctx)
	assert.NoError(mode.Do(ctx, &client, event))
}
