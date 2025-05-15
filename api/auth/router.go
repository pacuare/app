package auth

import (
	"net/http"

	"app.pacuare.dev/templates"
	"github.com/a-h/templ"
)

func Mount() {
	http.Handle("GET /auth/login", templ.Handler(templates.Login()))
}
