package www

import (
	"github.com/monkeydioude/moon/template"
)

// GetHome display home
func GetHome() ([]byte, int, error) {
	tmp := template.NewEngine("www/layout.html")

	tmp.
		BindTemplate("www/home.html", "Content").
		BindValues(template.HTML{"Title": "Salut Naomie"}).
		BindTemplate("www/list.html", "List").
		BindValue("bonjour", "ca va bien ?")

	return tmp.Render(), 200, nil
}
