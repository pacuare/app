package auth

import (
	"net/http"

	"dash.pacuare.dev/templates"
	"github.com/a-h/templ"
)

func Mount() {
	http.Handle("GET /auth/login", templ.Handler(templates.Login()))
}
