package handlers

import (
	"auth/components"
	"auth/db/models"
	"auth/db/repositories"
	"auth/session"
	"net/http"
)

// TODO: this should actually display online users not all of the users
func HandleOnlineUsers(w http.ResponseWriter, r *http.Request) {
	currentUserVal, _ := session.GetSessionValue(r, "username")
	currentUsername, ok := currentUserVal.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	users, err := repositories.GetUsernames()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Filter out the current user
	filtered := make([]models.User, 0, len(users))
	for _, user := range users {
		if user.Username != currentUsername {
			filtered = append(filtered, user)
		}
	}

	components.UserList(filtered).Render(r.Context(), w)
}
