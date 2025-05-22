package db

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type key int

const (
	keyDB key = iota
)

func DB(ctx context.Context) *pgxpool.Pool {
	return ctx.Value(keyDB).(*pgxpool.Pool)
}

func QueryOne[T any](ctx context.Context, query string, args ...any) (T, error) {
	var result T

	err := DB(ctx).
		QueryRow(context.Background(), query, args...).
		Scan(&result)

	return result, err
}

func ExportDB(conn *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keyDB, conn)))
		})
	}
}
