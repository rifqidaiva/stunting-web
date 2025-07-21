package handler

import (
	"html/template"
	"net/http"
	"path"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// Report handles the report page and serves the report template.
// Only allows GET requests; otherwise responds with 405 Method Not Allowed.
func Report(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var reportTemplate string = path.Join("web", "template", "report.html")
	var _head string = path.Join("web", "components", "_head.html")
	var _navbar string = path.Join("web", "components", "_navbar.html")
	var _footer string = path.Join("web", "components", "_footer.html")

	data := map[string]any{
		"document": map[string]any{
			"title": "Stunting Report",
			"meta": map[string]any{
				"description": "Report on stunting cases in Cirebon",
				"keywords":    "stunting, report, health, Cirebon",
			},
		},
	}

	template := template.Must(template.ParseFiles(reportTemplate, _head, _navbar, _footer))

	err := template.ExecuteTemplate(w, "report", data)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
