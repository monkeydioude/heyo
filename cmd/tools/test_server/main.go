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
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/monkeydioude/heyo/pkg/tiger/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const DefaultPort = "8022"

var ErrFailedToStartListener = errors.New("failed to start listener")
var ErrNewServerFromCreds = errors.New("error using TLS credentials")
var ErrServerServe = errors.New("failed to serve")

func boot() (string, *grpc.Server, *http.ServeMux) {
	port := DefaultPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	// listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	// assert.NoError(err, ErrFailedToStartListener)

	// TLS credentials for the backend connection
	creds, err := credentials.NewServerTLSFromFile("./certs/cert.pem", "./certs/key.pem")
	assert.NoError(err, ErrNewServerFromCreds)
	// Create the HTTP handler
	mux := http.NewServeMux()

	// REST API handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		log.Println("[INFO] received request on ", r.Method, r.URL)
		w.Write([]byte("Handled REST API request"))
	})
	return port, grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnknownServiceHandler((&CatchAllProxy{}).Proxy),
	), mux
}

func grpcHandler(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	port, grpcServer, mux := boot()
	assert.NotEmpty(port)
	assert.NotNil(mux)
	assert.NotNil(grpcServer)
	// Start the HTTP/2 server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: grpcHandler(grpcServer, mux),
	}
	log.Println("[INFO] Starting reverse proxy on", port)
	assert.NoError(server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"), ErrServerServe)
}
