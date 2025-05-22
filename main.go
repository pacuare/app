package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"app.pacuare.dev/api"
	"app.pacuare.dev/shared/auth"
	"app.pacuare.dev/shared/db"
	"app.pacuare.dev/templates"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Info("Connecting to database")
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	log.Info("Connected")
	defer conn.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(db.ExportDB(conn))
	r.Use(auth.Middleware)

	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.Mount("/api", api.Router())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		returnedApiKey := r.URL.Query().Get("key")
		email := auth.GetUser(r.Context())

		if email == nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		fullAccess, err := db.QueryOne[bool](r.Context(), "select fullAccess from AuthorizedUsers where email=$1", email)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		if !fullAccess {
			if databaseExists, err :=
				db.QueryOne[bool](r.Context(), "select count(*)>0 from pg_catalog.pg_database where datname = GetUserDatabase($1)", email); !databaseExists || err != nil {
				http.Redirect(w, r, "/createdb", http.StatusSeeOther)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}

		templates.Index(*email, fullAccess, returnedApiKey).Render(r.Context(), w)
	})

	http.HandleFunc("GET /createdb", func(w http.ResponseWriter, r *http.Request) {
		email := auth.GetUser(r.Context())

		if email == nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		if databaseExists, _ :=
			db.QueryOne[bool](r.Context(), "select count(*)>0 from pg_catalog.pg_database where datname = GetUserDatabase($1)", email); databaseExists {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templates.CreateDB(*email).Render(r.Context(), w)
		_, err = db.DB(r.Context()).Exec(context.Background(), fmt.Sprintf("create database %s", auth.GetUserDatabase(*email)))

		if err != nil {
			log.Error(err)
		}

		db, err := db.QueryOne[string](r.Context(), "select InitUserDatabase($1, $2, $3)", os.Getenv("DATABASE_URL_BASE"), os.Getenv("DATABASE_DATA"), *email)

		if err != nil {
			log.Error(err)
		}

		log.Infof("Created database %s", db)

		w.Write([]byte("<meta http-equiv=\"Refresh\" content=\"0;url=/\">"))
	})

	log.Info("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
