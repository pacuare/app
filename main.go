package main

import (
	"net/http"

	"app.pacuare.dev/api"
	"app.pacuare.dev/templates"
)

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	api.Mount()
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("authstatus")

		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		} else {
			templates.Index().Render(r.Context(), w)
		}
	})

	http.ListenAndServe(":8080", nil)
}
