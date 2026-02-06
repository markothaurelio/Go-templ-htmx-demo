package handlers

import (
	"net/http"

	"[repo]/news_article_app/templates"
	"github.com/a-h/templ"
)

func GraphHandler(cat string, graphs map[string]func() templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Collect all graph components
		var graphComponents []templ.Component
		for graphName, graphFunc := range graphs {
			graphComponents = append(graphComponents, templates.GraphComponent(graphName, graphFunc()))
		}

		// Render the layout with multiple graphs
		templates.Layout_Embedded(templates.RenderMultipleGraphs(graphComponents), "").Render(r.Context(), w)
	}
}
