package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/monkeydioude/moon"
	"github.com/monkeydioude/schampionne/bin/morab/www"
	"github.com/monkeydioude/tools"
)

const port = 6363

func getStaticResource(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	cssPath := r.Matches[0]

	css, err := ioutil.ReadFile(cssPath)

	if err != nil {
		return tools.Response404(err)
	}

	return tools.Response200(css)
}

func getHome(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	return www.GetHome()
}

func postHome(r *moon.Request, c *moon.Configuration) ([]byte, int, error) {
	return tools.Response200([]byte("ok"))
}

func isBrokerHealthy() error {
	res, err := tools.SendSimpleGetRequest(nil, nil, "http://localhost:19393/healthcheck")

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Healthcheck returned %d status code", res.StatusCode)
	}

	return nil
}

func main() {
	if err := isBrokerHealthy(); err != nil {
		log.Printf("[WARN] Broker could not be reached. Reason: %s", err)
	}
	log.Println("[INFO] Starting Morab")

	handler := moon.NewHandler(nil)
	handler.WithHeader("Access-Control-Allow-Origin", "*")

	handler.Routes.AddPost("^$", postHome)
	handler.Routes.AddGet("^$", getHome)
	handler.Routes.AddGet("^.+?\\.(css|js)$", getStaticResource)

	err := moon.ServerRun(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Printf("[ERR ] Server crashed. Reason: %s", err)
	}
}
