package handlers

import (
	"auth/components"
	"auth/db/models"
	"auth/db/repositories"
	"auth/services"
	"auth/session"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

// getting open channels (like general)
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

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	channelIDStr := chi.URLParam(r, "channel_id")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	msgs, err := repositories.GetMessagesFromChannelID(uint(channelID))
	if err != nil {
		http.Error(w, "Failed to load messages", http.StatusInternalServerError)
		return
	}

	components.Messages(msgs).Render(r.Context(), w)
}

func HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	channelIDStr := chi.URLParam(r, "channel_id")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	currentUser, _ := session.GetSessionValue(r, "username")
	username, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	user, err := repositories.FindUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	content := r.FormValue("content")
	err = services.CreateMessage(int(user.ID), int(channelID), content)
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

func ServeCommunication(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)
	currentUser, _ := session.GetSessionValue(r, "username")

	_, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	channelIDStr := chi.URLParam(r, "channel_id")
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}

	page := components.Base(
		"/chat",
		isAuthenticated,
		components.Communication(uint(channelID)),
	)
	templ.Handler(page).ServeHTTP(w, r)
}

func ServeChannelExplore(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)
	currentUser, _ := session.GetSessionValue(r, "username")

	username, ok := currentUser.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	page := components.Base("/channels/explore", isAuthenticated, components.ChatPage(username))
	templ.Handler(page).ServeHTTP(w, r)
}
