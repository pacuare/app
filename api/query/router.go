package query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"app.pacuare.dev/shared"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
)

func Mount() {
	http.HandleFunc("POST /api/query", func(w http.ResponseWriter, r *http.Request) {
		email, err := shared.GetUser(r)

		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(`{"error":"Not authorized"}`))
			return
		}

		language := r.URL.Query().Get("language")
		queryBuf := new(strings.Builder)
		_, err = io.Copy(queryBuf, r.Body)

		w.Header().Add("Content-Type", "application/json")

		if err != nil {
			log.Errorf("Error getting query: %e", err)

			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Internal server error"}`))
			return
		}

		query := queryBuf.String()

		if language == "python" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error":"Python is not yet supported"}`))
			return
		}

		// e.g. keeper@farthergate.com -> user_keeper__farthergate_com
		var conn *pgx.Conn
		var dbName string
		if fullAccess, err := shared.QueryOne[bool]("select fullAccess from AuthorizedUsers where email=$1", email); err != nil {
			log.Errorf("Error querying access level: %e", err)

			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Internal server error"}`))
			return
		} else if fullAccess {
			dbName = "pacuare_data"
		} else {
			dbName = shared.GetUserDatabase(*email)

		}

		dbUrl := fmt.Sprintf("%s/%s", os.Getenv("DATABASE_URL_BASE"), dbName)
		conn, err = pgx.Connect(r.Context(), dbUrl)

		if err != nil {
			log.Errorf("Error opening database: %e", err)

			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Internal server error"}`))
			return
		}

		defer conn.Close(r.Context())

		res, err := conn.Query(r.Context(), query)
		if err != nil {
			log.Errorf("Error running query: %e", err)

			w.WriteHeader(500)
			jsonError, _ := json.Marshal(map[string]string{"error": err.Error()})
			w.Write(jsonError)
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("["))

		firstRow := true

		for res.Next() {
			if firstRow {
				firstRow = false
			} else {
				w.Write([]byte(","))
			}
			values, _ := pgx.RowToMap(res)
			marshalled, _ := json.Marshal(values)
			w.Write(marshalled)
		}
		w.Write([]byte("]"))
	})
}
