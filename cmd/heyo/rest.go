package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/monkeydioude/heyo/internal/consts"
	"github.com/oklog/run"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"health\": \"OK\"}"))
}

func restRouting(mux *http.ServeMux) {
	mux.HandleFunc("/heyo/healthcheck", healthcheck)
}

func restServer(runGroup *run.Group) {
	mux := http.NewServeMux()
	restRouting(mux)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", consts.RestPort),
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}
	runGroup.Add(func() error {
		log.Println("[INFO] API starting on", consts.RestPort)
		return server.ListenAndServe()
	}, func(_ error) {
		log.Println("[INFO] closing API server")
		if err := server.Close(); err != nil {
			log.Println("[ERR ] ailed to stop web server", "err", err)
		}
	})
}
