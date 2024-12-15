package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monkeydioude/heyo/internal/service/boot"
	"github.com/monkeydioude/heyo/internal/service/client"
	"github.com/monkeydioude/heyo/internal/service/state"
	"github.com/monkeydioude/heyo/pkg/rpc"
	"github.com/monkeydioude/heyo/pkg/tiger/assert"
)

type Explorer struct {
	State  state.State
	ctx    context.Context
	client client.Client
}

func NewExplorer(ctx context.Context) Explorer {
	return Explorer{
		State:  state.Idle(),
		ctx:    ctx,
		client: boot.BootClient(ctx),
	}
}

func (e *Explorer) setupReceiverEvents() {
	e.ctx = e.client.GetCtx()
	receiverEvents, _ := readInput(e.ctx, "Events to subscribe to (comma separated)?")
	for _, event := range strings.Split(receiverEvents, ",") {
		event = strings.Trim(event, " \n")
		assert.NotEmpty(event)
		e.client.Listen(event, func(msg *rpc.Message) error {
			if e.State == state.STATE_IDLE && msg.ClientId != e.client.Uuid {
				fmt.Println("")
			}
			fmt.Printf("Message (%s) < %s\n", msg.ClientId, msg.Data)
			if e.State == state.STATE_IDLE && msg.ClientId != e.client.Uuid {
				fmt.Printf("Message (%s) > ", e.client.Uuid)
			}
			return nil
		})
	}
}

func (e *Explorer) setupName() {
	e.ctx = e.client.GetCtx()
	name, _ := readInput(e.ctx, "Your name?")
	name = strings.Trim(name, " \n")
	assert.NotEmpty(name)
	e.client.Uuid = name
}

func (e *Explorer) sendMessage(in *rpc.Message) error {
	e.State.Busy()
	if in == nil {
		return ErrInMessageWasNil
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	in.ClientId = e.client.Uuid
	ack, err := e.client.RpcClient.Enqueue(ctx, in)
	if err != nil {
		return fmt.Errorf("%w: %w", err, ErrEnqueuingMessage)
	}
	if ack.Code != rpc.AckCode_OK {
		return fmt.Errorf("%w: got error code %d", ErrAckNotOk, ack.Code)
	}
	e.State.Idle()
	return nil
}

func (e *Explorer) GetCtx() context.Context {
	return e.ctx
}
