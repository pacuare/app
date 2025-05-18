package main

import (
	"context"
	"net/http"
	"os"

	"app.pacuare.dev/api"
	"app.pacuare.dev/shared"
	"app.pacuare.dev/templates"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	log.Info("Connecting to database")
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	log.Info("Connected")
	shared.DB = conn
	defer conn.Close(context.Background())

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	api.Mount()
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		email, err := shared.GetUser(r)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		fullAccess, err := shared.QueryOne[bool]("select fullAccess from AuthorizedUsers where email=$1", email)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		}

		templates.Index(*email, fullAccess).Render(r.Context(), w)
	})

	log.Info("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
