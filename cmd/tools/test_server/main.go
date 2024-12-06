// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"time"
// )

// func restServer() error {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Printf("io.ReadAll err: %s\n", err)
// 		}
// 		fmt.Println("incoming request", r.Method, r.URL.Path, string(body))
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("ok"))
// 	})
// 	server := &http.Server{
// 		Addr:              ":8022",
// 		ReadTimeout:       3 * time.Second,
// 		WriteTimeout:      3 * time.Second,
// 		IdleTimeout:       30 * time.Second,
// 		ReadHeaderTimeout: 2 * time.Second,
// 		Handler:           mux,
// 	}
// 	log.Println("API starting on 8022")
// 	return server.ListenAndServe()
// }

// func main() {
// 	if err := restServer(); err != nil {
// 		panic(err)
// 	}
// }

package main

import (
	"crypto/tls"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// func grpcCatchAllHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if the incoming request is a gRPC request
// 	if !(r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc") {
// 		http.Error(w, "Not Found", http.StatusNotFound)
// 	}

// 	fmt.Println("Received a gRPC request")
// 	fmt.Printf("Method: %s, Path: %s\n", r.Method, r.URL.Path)

// 	// Read the raw body (which contains gRPC data)
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Printf("Body: %x\n", body) // Print the raw body in hexadecimal for debugging

// 	// Send a simple gRPC response (gRPC requires a specific format)
// 	w.Header().Set("Content-Type", "application/grpc")
// 	w.WriteHeader(http.StatusOK)

// 	// gRPC responses must include a status trailer
// 	w.Write([]byte{})                  // This could be serialized protobuf data if you wanted
// 	w.Header().Set("Grpc-Status", "0") // "0" means OK in gRPC
// }

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
		forwardStreamClient(stream, clientStream, errCh)
	}()
	go func() {
		forwardStreamServer(clientStream, stream, errCh)
	}()
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}
	return nil
}

func forwardStreamServer(src grpc.ClientStream, dst grpc.ServerStream, errCh chan error) {
	for {
		// Use a generic interface for receiving messages
		// var m any // This ensures the data is treated as a Protobuf message
		// var m :=
		m := new(proto.Message)
		err := src.RecvMsg(m)
		if err != nil {
			errCh <- err
			return
		}

		// Forward the received message to the destination
		err = dst.SendMsg(m)
		if err != nil {
			errCh <- err
			return
		}
	}
}
func forwardStreamClient(src grpc.ServerStream, dst grpc.ClientStream, errCh chan error) {
	for {
		// Use a generic interface for receiving messages
		// var m any // This ensures the data is treated as a Protobuf message
		// var m :=
		var m proto.Message
		err := src.RecvMsg(&m)
		if err != nil {
			errCh <- err
			return
		}

		// Forward the received message to the destination
		err = dst.SendMsg(m)
		if err != nil {
			errCh <- err
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8022")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	// TLS credentials for the backend connection
	creds, err := credentials.NewServerTLSFromFile("./certs/cert.pem", "./certs/key.pem")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnknownServiceHandler((&CatchAllProxy{}).Proxy),
	)

	log.Println("Starting reverse proxy on :8022")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
