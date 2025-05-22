package handlers

import (
	"auth/components"
	"auth/db/repositories"
	"auth/session"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	username, err := session.GetSessionValue(r, "username")
	if err != nil || username == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := username.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	isAuthenticated := session.IsAuthenticated(r)

	found_user, err := repositories.FindUserByUsername(user)
	if err != nil {
		log.Printf("Failed to find user with username %s: %v", username, err)
		return
	}

	page := components.Base("/dashboard", isAuthenticated, components.Dashboard(user, found_user.PGPKey))
	templ.Handler(page).ServeHTTP(w, r)
}
