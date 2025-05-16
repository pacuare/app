package api

import (
	"fmt"
	"net/http"

	"app.pacuare.dev/api/auth"
	"github.com/jackc/pgx/v5"
)

func Mount(conn *pgx.Conn) {
	auth.Mount(conn)

	http.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, `{"ok":"true"}`)
	})
}
