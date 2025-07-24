package main

import (
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/api"
	"github.com/rifqidaiva/stunting-web/internal/object"
)

func main() {
	// http.Handle("/static/",
	// 	http.StripPrefix("/static/",
	// 		http.FileServer(http.Dir("assets"))))

	// http.HandleFunc("/", handler.Index)
	// http.HandleFunc("/login", handler.Login)
	// http.HandleFunc("/report", handler.Report)
	// http.HandleFunc("/admin", handler.Admin)

	http.HandleFunc("/api/admin/insert", api.AdminInsert)
	http.HandleFunc("/api/admin/get", api.AdminGet)
	http.HandleFunc("/api/admin/get/geojson", api.AdminGetGeoJson)
	http.HandleFunc("/api/admin/update", api.AdminUpdate)
	http.HandleFunc("/api/admin/delete", api.AdminDelete)

	// API test endpoint
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		response := object.NewResponse(http.StatusOK, "Test API is working", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
