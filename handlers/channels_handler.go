package handlers

import (
	"auth/components"
	"auth/db/repositories"
	"auth/session"
	"log"
	"net/http"
)

// TODO: Get No Direct Channels
func GetOpenChannels(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := session.GetSessionValue(r, "username")

	username, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	user, err := repositories.FindUserByUsername(username)
	if err != nil {
		log.Printf("Failed to find user with username %s: %v", username, err)
		return
	}

	channels, err := repositories.FetchOpenChannels()
	if err != nil {
		http.Error(w, "Failed to fetch open channels", http.StatusInternalServerError)
		return
	}
	components.OpenChannelsList(channels, user.ID).Render(r.Context(), w)
}
