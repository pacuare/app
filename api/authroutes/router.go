package authroutes

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app.pacuare.dev/api/authroutes/mailer"
	"app.pacuare.dev/shared/db"
	"app.pacuare.dev/shared/enc"
	"app.pacuare.dev/templates"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		failedEmail := r.URL.Query().Get("failed-email")
		var pFailedEmail *string = nil

		if r.URL.Query().Has("failed-email") {
			pFailedEmail = &failedEmail
		}

		templates.Login(pFailedEmail).Render(r.Context(), w)
	})

	r.Get("/verify", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")

		hasUser, err := db.QueryOne[bool](r.Context(), "select (count(*) > 0) from AuthorizedUsers where email=$1", email)

		if err != nil {
			log.Error(err)
		}

		if !hasUser {
			http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s", email), http.StatusSeeOther)
		} else {
			_, err := mailer.SendConfirmation(db.DB(r.Context()), email)

			if err != nil {
				log.Errorf("Error creating verification code for %s", email)
				http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s", email), http.StatusSeeOther)
				return
			}

			templates.Verify(email).Render(r.Context(), w)
		}
	})

	http.HandleFunc("POST /auth/verify", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")

		err := r.ParseForm()
		if err != nil {
			log.Error(err)
			http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s&login-failed", email), http.StatusSeeOther)
			return
		}

		inputCode := r.Form.Get("otp")

		log.Infof("%s %s", email, inputCode)

		code, err := db.QueryOne[string](r.Context(), "select code from LoginCodes where email=$1", email)

		if err != nil {
			log.Error(err)
			http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s&login-failed", email), http.StatusSeeOther)
			return
		}

		if strings.EqualFold(code, inputCode) {
			enc, err := enc.Encrypt([]byte(email))
			if err != nil {
				log.Error(err)
				http.Redirect(w, r, fmt.Sprintf("/auth/login?failed-email=%s&login-failed", email), http.StatusSeeOther)
				return
			}

			c := &http.Cookie{Name: "AuthStatus", Value: hex.EncodeToString(enc), MaxAge: 259200, Path: "/", HttpOnly: true}
			http.SetCookie(w, c)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "AuthStatus", Path: "/", Expires: time.Unix(0, 0)})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	return r
}
