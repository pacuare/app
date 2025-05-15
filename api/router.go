package api

import (
	"app.pacuare.dev/api/auth"
	"github.com/jackc/pgx/v5"
)

func Mount(conn *pgx.Conn) {
	auth.Mount(conn)
}
