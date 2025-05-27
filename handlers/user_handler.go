package handlers

import (
	"auth/components"
	"auth/db/repositories"
	"net/http"
)

// TODO: this should actually display online users not all of the users
func HandleOnlineUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repositories.GetUsernames()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	components.UserList(users).Render(r.Context(), w)
}
