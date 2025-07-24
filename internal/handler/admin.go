package handler

import (
	"html/template"
	"net/http"
	"path"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// Admin handles the admin page and serves the admin template.
// Only allows GET requests; otherwise responds with 405 Method Not Allowed.
func Admin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var edit string = path.Join("web", "template", "admin.html")
	var _head string = path.Join("web", "components", "_head.html")
	var _navbar string = path.Join("web", "components", "_navbar.html")
	var _footer string = path.Join("web", "components", "_footer.html")

	var _addSuffererContent string = path.Join("web", "components", "admin", "_addSuffererContent.html")
	var _dataContent string = path.Join("web", "components", "admin", "_dataContent.html")
	var _statisticsContent string = path.Join("web", "components", "admin", "_statisticsContent.html")

	data := map[string]any{
		"document": map[string]any{
			"page":  "admin",
			"title": "Stunting Kota Cirebon",
			"meta": map[string]any{
				"description": "Edit Stunting Data - Stunting Kota Cirebon",
				"keywords":    "stunting, kesehatan, anak, Cirebon",
			},
		},
	}

	template := template.Must(template.ParseFiles(edit, _head, _navbar, _footer, _addSuffererContent, _dataContent, _statisticsContent))

	err := template.ExecuteTemplate(w, "admin", data)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
