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
	"net"
	"os"

	"github.com/monkeydioude/heyo/pkg/tiger/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const DefaultPort = "8022"

var ErrFailedToStartListener = errors.New("failed to start listener")
var ErrNewServerFromCreds = errors.New("error using TLS credentials")
var ErrServerServe = errors.New("failed to serve")

func boot() (net.Listener, *grpc.Server) {
	port := DefaultPort
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	assert.NoError(err, ErrFailedToStartListener)

	// TLS credentials for the backend connection
	creds, err := credentials.NewServerTLSFromFile("./certs/cert.pem", "./certs/key.pem")
	assert.NoError(err, ErrNewServerFromCreds)
	return listener, grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnknownServiceHandler((&CatchAllProxy{}).Proxy),
	)
}

func main() {
	listener, server := boot()
	assert.NotNil(listener)
	assert.NotNil(server)
	log.Println("[INFO] Starting reverse proxy on", listener.Addr())
	assert.NoError(server.Serve(listener), ErrServerServe)
}
