package api

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"app.pacuare.dev/api/authroutes"
	"app.pacuare.dev/api/query"
	"app.pacuare.dev/shared/auth"
	"app.pacuare.dev/shared/db"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed openapi.yml
var apiSpec string

func Router() chi.Router {
	r := chi.NewRouter()

	r.Mount("/auth", authroutes.Router())
	r.Mount("/query", query.Router())

	r.HandleFunc("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodGet || r.Method == http.MethodOptions) {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/yaml")
		w.WriteHeader(200)
		fmt.Fprint(w, apiSpec)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, `{"ok":"true"}`)
	})

	r.Use(auth.RequireAuth)

	r.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		email := auth.GetUser(r.Context())

		db, err := db.QueryOne[string](r.Context(), "select InitUserDatabase($1, $2, $3)", os.Getenv("DATABASE_URL_BASE"), os.Getenv("DATABASE_DATA"), *email)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(db))
	})

	r.Post("/recreate", func(w http.ResponseWriter, r *http.Request) {
		email := auth.GetUser(r.Context())

		_, err := r.Context().Value("db").(*pgxpool.Pool).Exec(context.Background(), fmt.Sprintf("drop database %s", auth.GetUserDatabase(*email)))

		if err != nil {
			log.Error(err)
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Deleted successfully"))
	})

	r.Post("/key", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/?settings", http.StatusSeeOther)
			return
		}

		email := auth.GetUser(r.Context())
		description := r.FormValue("description")

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		key, err := db.QueryOne[string](r.Context(), "insert into APIKeys (email, description) values ($1, $2) returning key", email, description)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/?settings", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/?settings&key=%s", key), http.StatusSeeOther)
	})

	r.Post("/key/delete", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/?settings", http.StatusSeeOther)
			return
		}

		email := auth.GetUser(r.Context())
		id := r.FormValue("id")

		_, err = db.DB(r.Context()).Exec(context.Background(), "delete from APIKeys where id = $1 and email = $2", id, email)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, "/?settings", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/?settings", http.StatusSeeOther)
	})

	return r
}
