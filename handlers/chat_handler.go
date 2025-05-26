package handlers

import (
	"auth/components"
	"auth/db/models"
	"auth/db/repositories"
	"auth/services"
	"auth/session"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func ServeChatPage(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)
	currentUser, _ := session.GetSessionValue(r, "username")

	user, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	page := components.Base("/chat", isAuthenticated, components.ChatPage(user))
	// log.Println("Route has been hit") - it is printing
	templ.Handler(page).ServeHTTP(w, r)
}

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Route has been hit")
	msgs, err := repositories.GetRecentMessages()
	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	components.Messages(msgs).Render(r.Context(), w)
}

func HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Route has been hit")
	currentUser, _ := session.GetSessionValue(r, "username")

	username, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	log.Printf("I should get username: %s", username)

	user, err := repositories.FindUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	content := r.FormValue("content")

	log.Printf("I should get content: %s", content)

	err = services.CreateMessage(int(user.ID), content)
	if err != nil {
		http.Error(w, "Creating new Message has failed", http.StatusInternalServerError)
		return
	}

	msg := models.MessageWithUser{
		Username:  user.Username,
		Content:   content,
		CreatedAt: time.Now(),
	}
	components.Messages([]models.MessageWithUser{msg}).Render(r.Context(), w)
}
