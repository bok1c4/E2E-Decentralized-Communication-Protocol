package handlers

import (
	"auth/components"
	"auth/db/repositories"
	"auth/session"
	"auth/util"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func ServeGenPGP(w http.ResponseWriter, r *http.Request) {
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

	page := components.Base("/getpgp", isAuthenticated, components.GetPGP(user))
	templ.Handler(page).ServeHTTP(w, r)
}

func HandleGenPGP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

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

	pgpKey := r.FormValue("pgpKey")
	if pgpKey == "" {
		http.Error(w, "PGP Key is required", http.StatusBadRequest)
		return
	}

	if err := util.IsValidPublicKey(pgpKey); err != nil {
		http.Error(w, "Invalid PGP key: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = repositories.InsertPGPKey(user, pgpKey)
	if err != nil {
		log.Printf("Failed to insert PGP key for %s: %v", user, err)
		http.Error(w, "Failed to save PGP key", http.StatusInternalServerError)
		return
	}

	log.Printf("PGP key successfully inserted for user %s", user)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
