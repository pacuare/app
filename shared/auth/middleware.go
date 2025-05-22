package auth

import (
	"context"
	"encoding/hex"
	"net/http"
	"strings"

	"app.pacuare.dev/shared/db"
	"app.pacuare.dev/shared/enc"
)

type key int

const (
	authKey key = iota
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authStatus, err := r.Cookie("AuthStatus")

		if err == nil {
			authBytes, err := hex.DecodeString(authStatus.Value)
			if err != nil {

				return
			}

			email, err := enc.Decrypt(authBytes)
			if err != nil {
				return
			}

			emailStr := string(email)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authKey, &emailStr)))
		} else {
			authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			emailStr, err := db.QueryOne[string](r.Context(), "select email from APIKeys where key = $1", authHeader)
			if err != nil {
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authKey, &emailStr)))
		}
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUser(r.Context()) == nil {
			w.WriteHeader(401)
			w.Write([]byte("Not authorized"))
			return
		}
	})
}

func GetUser(ctx context.Context) *string {
	return ctx.Value(authKey).(*string)
}
