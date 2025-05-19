package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"app.pacuare.dev/api/auth"
	"app.pacuare.dev/api/query"
	"app.pacuare.dev/shared"
	"github.com/charmbracelet/log"
)

func Mount() {
	auth.Mount()
	query.Mount()

	http.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, `{"ok":"true"}`)
	})

	http.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		email, err := shared.GetUser(r)

		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Not authorized"))
			return
		}

		db, err := shared.QueryOne[string]("select InitUserDatabase($1, $2, $3)", os.Getenv("DATABASE_URL_BASE"), os.Getenv("DATABASE_DATA"), *email)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(db))
	})

	http.HandleFunc("POST /api/recreate", func(w http.ResponseWriter, r *http.Request) {
		email, err := shared.GetUser(r)

		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Not authorized"))
			return
		}

		_, err = shared.DB.Exec(context.Background(), fmt.Sprintf("drop database %s", shared.GetUserDatabase(*email)))

		if err != nil {
			log.Error(err)
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("Deleted successfully"))
	})
}
