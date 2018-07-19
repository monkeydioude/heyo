package www

import (
	"github.com/monkeydioude/moon/template"
)

// GetHome display home
func GetHome() ([]byte, int, error) {
	tmp := template.NewEngine("www/layout.html")

	tmp.
		BindTemplate("www/home.html", "Content").
		AddValue("Title", "Salut").
		// AddValue("Body", "C'est cool")
		BindTemplate("www/list.html", "List").
		AddValue("bonjour", "bien ?")

	return tmp.Render(), 200, nil
}
