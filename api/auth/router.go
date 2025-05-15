package auth

import (
	"context"
	"fmt"
	"net/http"

	"app.pacuare.dev/templates"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func Mount(conn *pgx.Conn) {
	http.HandleFunc("GET /auth/login", func(w http.ResponseWriter, r *http.Request) {
		failedEmail := r.URL.Query().Get("failed-email")
		var pFailedEmail *string = nil

		if r.URL.Query().Has("failed-email") {
			pFailedEmail = &failedEmail
		}

		templates.Login(pFailedEmail).Render(r.Context(), w)
	})

	http.HandleFunc("GET /auth/verify", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")

		var hasUser bool
		err := conn.
			QueryRow(context.Background(), "select (count(*) > 0) from AuthorizedUsers where email=$1", email).
			Scan(&hasUser)

		if err != nil {
			log.Error(err)
		}

		if !hasUser {
			http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s", email), http.StatusSeeOther)
		} else {
			templates.Verify(email).Render(r.Context(), w)
		}
	})
}
