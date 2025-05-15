package api

import "app.pacuare.dev/api/auth"

func Mount() {
	auth.Mount()
}
