package main

import (
	"context"
	"fmt"
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
			return
		}

		if !fullAccess {
			if databaseExists, err :=
				shared.QueryOne[bool]("select count(*)>0 from pg_catalog.pg_database where datname = GetUserDatabase($1)", email); !databaseExists || err != nil {
				http.Redirect(w, r, "/createdb", http.StatusSeeOther)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}

		templates.Index(*email, fullAccess).Render(r.Context(), w)
	})

	http.HandleFunc("GET /createdb", func(w http.ResponseWriter, r *http.Request) {
		email, err := shared.GetUser(r)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		if databaseExists, _ :=
			shared.QueryOne[bool]("select count(*)>0 from pg_catalog.pg_database where datname = GetUserDatabase($1)", email); databaseExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templates.CreateDB(*email).Render(r.Context(), w)
		_, err = shared.DB.Exec(context.Background(), fmt.Sprintf("create database %s", shared.GetUserDatabase(*email)))

		if err != nil {
			log.Error(err)
		}

		db, err := shared.QueryOne[string]("select InitUserDatabase($1, $2, $3)", os.Getenv("DATABASE_URL_BASE"), os.Getenv("DATABASE_DATA"), *email)

		if err != nil {
			log.Error(err)
		}

		log.Infof("Created database %s", db)

		w.Write([]byte("<meta http-equiv=\"Refresh\" content=\"0;url=/\">"))
	})

	log.Info("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
