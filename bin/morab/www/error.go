package www

import (
	"bytes"
	"html/template"
)

type derror struct {
	Code int
}

func err500() ([]byte, int, error) {
	t, _ := template.ParseFiles("www/err_code.html")
	b := &bytes.Buffer{}
	t.Execute(b, derror{500})
	return b.Bytes(), 500, nil
}
