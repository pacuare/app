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

func queryEndpointArgs(w http.ResponseWriter, r *http.Request, query string, params []any) {
	email, err := shared.GetUser(r)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Not authorized"}`))
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

func Mount() {
	http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		if !(r.Method == http.MethodPost || r.Method == http.MethodOptions) {
			w.WriteHeader(405)
			fmt.Fprint(w, "Method not allowed")
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
}
