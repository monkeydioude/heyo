package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/monkeydioude/heyo/pkg/rpc"
)

func readInput(ctx context.Context, msg string) (string, error) {
	result := make(chan string)
	errChan := make(chan error)
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s > ", msg)

	// Run the read operation in a separate goroutine
	go func() {
		line, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
		} else {
			result <- line
		}
	}()

	select {
	case <-ctx.Done(): // Context canceled
		return "", fmt.Errorf("%w: %w", ErrUpstreamClosed, ctx.Err())
	case err := <-errChan: // Read error
		return "", err
	case line := <-result: // Successfully read a line
		return line, nil
	}
}

func buildMessage(ctx context.Context, ClientId string) (rpc.Message, error) {
	input, err := readInput(ctx, fmt.Sprintf("Message (%s)", ClientId))
	if err != nil {
		return rpc.Message{}, fmt.Errorf("%w: %w", ErrReadingInput, err)
	}
	parts := strings.SplitN(input, " ", 2)

	if len(parts) != 2 {
		return rpc.Message{}, fmt.Errorf("%w: %w", ErrMessageWasMalformed, ErrMessageMinimum2Words)
	}
	if parts[0][0] != '@' {
		return rpc.Message{}, fmt.Errorf("%w: %w", ErrMessageWasMalformed, ErrLeadingAtMissing)
	}
	return rpc.Message{
		Event:     strings.Trim(parts[0], " @\n"),
		Data:      strings.Trim(parts[1], " \n"),
		MessageId: uuid.NewString(),
	}, nil
}
