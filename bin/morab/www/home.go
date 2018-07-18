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
	// t, err := template.ParseFiles("www/home.html")

	// if err != nil {
	// 	log.Printf("[ERR ] Could not parse file. Reason: %s", err)
	// 	return err500()
	// }
	// b := &bytes.Buffer{}
	// t.Execute(b, map[string]string{"test": "bonjour"})

	// return b.Bytes(), 200, nil

	tmp := template.NewEngine("www/layout.html")
	tmp.WithChild("www/home.html", "Content")
	tmp.AddValue(&Home{
		Title: "Bonjour",
		Body:  "c'est cool",
	})
	return tmp.Render(), 200, nil
}
