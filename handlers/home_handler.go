package handlers

import (
	"auth/components"
	"auth/session"
	"net/http"

	"github.com/a-h/templ"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// check if user is authenticated
	// check if user have pgp key
	isAuthenticated := session.IsAuthenticated(r)

	page := components.Base("/", isAuthenticated, components.Hero(isAuthenticated))
	templ.Handler(page).ServeHTTP(w, r)
}
