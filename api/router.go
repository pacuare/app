package api

import "dash.pacuare.dev/api/auth"

func Mount() {
	auth.Mount()
}
