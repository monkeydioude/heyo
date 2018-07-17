package www

import (
	"bytes"
	"html/template"
	"log"
)

type home struct {
	Msg  string
	Code int
}

// GetHome display home
func GetHome() ([]byte, int, error) {
	t, err := template.ParseFiles("www/home.html")

	if err != nil {
		log.Printf("[ERR ] Could not parse file. Reason: %s", err)
		return err500()
	}
	b := &bytes.Buffer{}
	t.Execute(b, home{"wesh", 200})

	return b.Bytes(), 200, nil
}
