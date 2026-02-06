package handlers

import (
	"net/http"

	"[repo]/news_article_app/templates"
	"github.com/a-h/templ"
)

func MapHandler(cat string, maps map[string]map[string]func(path string) templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year := r.URL.Query().Get("year")

		// Render the layout with the correct map template
		templates.Layout_Embedded(templates.Maps(year, maps), "").Render(r.Context(), w)
	}
}
