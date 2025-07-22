package main

import (
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/api"
	"github.com/rifqidaiva/stunting-web/internal/handler"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/report", handler.Report)
	http.HandleFunc("/admin", handler.Admin)

	http.HandleFunc("/api/admin/insert", api.AdminInsert)
	http.HandleFunc("/api/admin/get", api.AdminGet)
	http.HandleFunc("/api/admin/get/geojson", api.AdminGetGeoJson)
	http.HandleFunc("/api/admin/update", api.AdminUpdate)
	http.HandleFunc("/api/admin/delete", api.AdminDelete)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
