package main

import (
	"fmt"
	"net/http"

	_ "github.com/rifqidaiva/stunting-web/docs" // Import for Swagger documentation
	"github.com/rifqidaiva/stunting-web/internal/object"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Stunting Web API
// @version 0.0.2
// @description API for managing stunting data
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func main() {
	// http.HandleFunc("/api/admin/get", api.AdminGet)

	// http.HandleFunc("/api/admin/get", api.AdminGet)
	// http.HandleFunc("/api/admin/update", api.AdminUpdate)

	// API test endpoint
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		response := object.NewResponse(http.StatusOK, "Test API is working", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/swagger/doc.json", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
