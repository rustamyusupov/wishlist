package main

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	if tmpl == nil {
		if tmpl == nil {
			tmpl = template.Must(tmpl.ParseGlob("views/*.html"))
		}
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Whishlist",
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}
