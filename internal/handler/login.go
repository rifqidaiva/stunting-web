package handler

import (
	"html/template"
	"net/http"
	"path"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// Login handles the login page and serves the login template.
// Only allows GET requests; otherwise responds with 405 Method Not Allowed.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var loginTemplate string = path.Join("web", "template", "login.html")
	var _head string = path.Join("web", "components", "_head.html")
	var _navbar string = path.Join("web", "components", "_navbar.html")
	var _footer string = path.Join("web", "components", "_footer.html")

	data := map[string]any{
		"document": map[string]any{
			"page": "login",
			"title": "Login - Stunting Kota Cirebon",
			"meta": map[string]any{
				"description": "Login to the Stunting Kota Cirebon system",
				"keywords":    "login, stunting, kesehatan, anak, Cirebon",
			},
		},
	}

	template := template.Must(template.ParseFiles(loginTemplate, _head, _navbar, _footer))

	err := template.ExecuteTemplate(w, "login", data)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
