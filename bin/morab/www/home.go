package www

import (
	"github.com/monkeydioude/moon/template"
)

type Home struct {
	Title string
	Body  string
}

// GetHome display home
func GetHome() ([]byte, int, error) {
	tmp := template.NewEngine("www/layout.html")
	tmp.
		WithChild("www/home.html", "Content").
		AddValue(&Home{
			Title: "Bonjour",
			// Body:  "c'est cool",
		}).
		WithChild("www/list.html", "List")
	return tmp.Render(), 200, nil
}
