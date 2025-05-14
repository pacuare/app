package main

import (
	"net/http"

	"github.com/a-h/templ"
	"pacuare.dev/dash/templates"
)

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/", templ.Handler(templates.Index()))

	http.ListenAndServe(":8080", nil)
}
