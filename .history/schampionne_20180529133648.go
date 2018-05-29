package main

import (
	"context"

	"github.com/monkeydioude/schampionne/grpc"
)

type server struct{}

func (s *server) Whisper(context.Context, *grpc.Rumor) (*grpc.Response, error) {
	return nil, nil
}
