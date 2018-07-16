package main

import (
	"fmt"
	"log"

	"github.com/monkeydioude/moon"
	"github.com/monkeydioude/tools"
)

const port = 19393

func getHome(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	return tools.Response200([]byte("ok"))
}

func postHome(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	return tools.Response200([]byte("ok"))
}

func isBrokerHealthy() error {
	res, err := tools.SendSimpleGetRequest(nil, nil, "http://localhost:6363/healthcheck")

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Healthcheck returned %d status code", res.StatusCode)
	}

	return nil
}

func main() {
	// use goroutine
	if err := isBrokerHealthy(); err != nil {
		log.Fatalf("[ERR ] Could not start Morab. Reason: %s", err)
	}
	log.Println("[INFO] Broker's up, starting Morab")

	handler := moon.NewHandler(nil)
	handler.WithHeader("Access-Control-Allow-Origin", "*")

	handler.Routes.AddPost("^$", postHome)
	handler.Routes.AddGet("^$", getHome)

	err := moon.ServerRun(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Printf("[ERR ] Server crashed. Reason: %s", err)
	}
}
