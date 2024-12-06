package main

import (
	"crypto/tls"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CatchAllProxy intercepts all incoming gRPC requests and forwards them to the backend
type CatchAllProxy struct{}

// Proxy implements the generic proxy functionality
func (p *CatchAllProxy) Proxy(srv interface{}, stream grpc.ServerStream) error {
	// Extract the full method name (e.g., /ServiceName/MethodName)
	methodName, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Errorf(codes.Unavailable, "Unable to get method name from stream")
	}
	log.Printf("Received request for method: %s", methodName)

	// Extract incoming metadata
	md, _ := metadata.FromIncomingContext(stream.Context())
	log.Printf("Incoming metadata: %v", md)

	// Backend server address
	backendAddress := "localhost:9393" // Replace with dynamic resolution if needed

	// // TLS credentials for the backend connection
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // Skip verification for testing; remove this in production
	})

	// creds, err := credentials.NewServerTLSFromFile("./certs/cert.pem", "./certs/key.pem")
	// if err != nil {
	// 	return status.Errorf(codes.Unavailable, "Unable to connect to backend: %v", err)
	// }
	// Dial the backend server using grpc.NewClient
	conn, err := grpc.NewClient(
		backendAddress,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return status.Errorf(codes.Unavailable, "Unable to connect to backend: %v", err)
	}
	defer conn.Close()

	// Forward the request to the backend
	clientStream, err := grpc.NewClientStream(
		stream.Context(),
		&grpc.StreamDesc{
			ServerStreams: true,
			ClientStreams: true,
		},
		conn,
		methodName,
	)
	if err != nil {
		return status.Errorf(codes.Unavailable, "Unable to create client stream: %v", err)
	}

	// Proxy all requests and responses between the client and backend
	errCh := make(chan error, 2)
	go func() {
		forwardStream(stream, clientStream, errCh)
	}()
	go func() {
		forwardStream(clientStream, stream, errCh)
	}()
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}
	return nil
}
