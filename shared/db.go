package shared

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func QueryOne[T any](query string, args ...any) (T, error) {
	var result T

	err := DB.
		QueryRow(context.Background(), query, args...).
		Scan(&result)

	return result, err
}
