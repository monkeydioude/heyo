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
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func grpcCatchAllHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the incoming request is a gRPC request
	if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
		fmt.Println("Received a gRPC request")
		fmt.Printf("Method: %s, Path: %s\n", r.Method, r.URL.Path)

		// Read the raw body (which contains gRPC data)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Body: %x\n", body) // Print the raw body in hexadecimal for debugging

		// Send a simple gRPC response (gRPC requires a specific format)
		w.Header().Set("Content-Type", "application/grpc")
		w.WriteHeader(http.StatusOK)

		// gRPC responses must include a status trailer
		w.Write([]byte{})                  // This could be serialized protobuf data if you wanted
		w.Header().Set("Grpc-Status", "0") // "0" means OK in gRPC
	} else {
		// If it's not a gRPC request, respond with a 404
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func grpcServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", grpcCatchAllHandler)

	server := &http.Server{
		Addr:              ":8022",
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}

	log.Println("gRPC server starting on port 8022")
	return server.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem") // gRPC typically requires TLS
}

func main() {
	if err := grpcServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
