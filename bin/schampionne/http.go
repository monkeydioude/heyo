package main

import (
	"log"

	"github.com/monkeydioude/moon"
	"github.com/monkeydioude/tools"
)

func healthcheck(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	return tools.Response200([]byte("ok"))
}

func httpServer() {
	handler := moon.NewHandler(nil)
	handler.Routes.AddGet("healthcheck$", healthcheck)

	if err := moon.ServerRun(httpPort, handler); err != nil {
		log.Printf("[ERR ] Server crashed. Reason: %s", err)
	}
}
