package query

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"app.pacuare.dev/shared/auth"
	"app.pacuare.dev/shared/db"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func queryEndpointArgs(w http.ResponseWriter, r *http.Request, query string, params []any) {
	email := auth.GetUser(r.Context())

	if email == nil {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"Not authorized"}`))
		return
	}

	// e.g. keeper@farthergate.com -> user_keeper__farthergate_com
	var conn *pgx.Conn
	var dbName string
	if fullAccess, err := db.QueryOne[bool](r.Context(), "select fullAccess from AuthorizedUsers where email=$1", email); err != nil {
		log.Errorf("Error querying access level: %e", err)

		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	} else if fullAccess {
		dbName = "pacuare_data"
	} else {
		dbName = auth.GetUserDatabase(*email)

	}

	dbUrl := fmt.Sprintf("%s/%s", os.Getenv("DATABASE_URL_BASE"), dbName)
	conn, err := pgx.Connect(r.Context(), dbUrl)

	if err != nil {
		log.Errorf("Error opening database: %e", err)

		w.WriteHeader(500)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}

	defer conn.Close(r.Context())

	res, err := conn.Query(r.Context(), query, params...)
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
}

func Router() chi.Router {
	r := chi.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodPost || r.Method == http.MethodOptions) {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			w.Write([]byte(`{"preflight": "ok"}`)) // send the preflight on its way before we hit logic
			return
		}

		var params []any
		var query string

		reqBuf := new(strings.Builder)

		_, err := io.Copy(reqBuf, r.Body)

		if err != nil {
			log.Errorf("Error getting query: %e", err)

			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Internal server error"}`))
			return
		}

		if r.Header.Get("Content-Type") == "application/json" {
			var jsonBody struct {
				Query  string `json:"query"`
				Params []any  `json:"params"`
			}

			err = json.Unmarshal([]byte(reqBuf.String()), &jsonBody)

			if err != nil {
				log.Error(err)
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte(`{"error":"Parse error"}`))
				return
			}
			params = jsonBody.Params
			query = jsonBody.Query
		} else {
			params = []any{}
			query = reqBuf.String()
		}

		queryEndpointArgs(w, r, query, params)
	})

	return r
}
