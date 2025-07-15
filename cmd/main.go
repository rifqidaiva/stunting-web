package main

import (
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/handler"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/edit", handler.Edit)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
