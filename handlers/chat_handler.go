package handlers

import (
	"auth/components"
	"auth/db/models"
	"auth/db/repositories"
	"auth/services"
	"auth/session"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

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

func ServeStartCommunication(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)

	currentUserVal, _ := session.GetSessionValue(r, "username")
	currentUsername, ok := currentUserVal.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	currentUser, err := repositories.FindUserByUsername(currentUsername)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	targetUsername := chi.URLParam(r, "username")
	if targetUsername == "" || targetUsername == currentUsername {
		http.NotFound(w, r)
		return
	}

	targetUser, err := repositories.FindUserByUsername(targetUsername)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	channel, err := repositories.FindChanBetweenTwoUsers(currentUser.ID, targetUser.ID)
	if err == nil && channel != nil {
		http.Redirect(w, r, "/channel/"+strconv.FormatUint(uint64(channel.ID), 10), http.StatusSeeOther)
		return
	}

	page := components.Base(
		"/chat/"+targetUsername,
		isAuthenticated,
		components.StartConvo(targetUsername),
	)
	templ.Handler(page).ServeHTTP(w, r)
}

func HandleChatInitSend(w http.ResponseWriter, r *http.Request) {
	currentUserVal, _ := session.GetSessionValue(r, "username")
	currentUsername, ok := currentUserVal.(string)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	sender, err := repositories.FindUserByUsername(currentUsername)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	targetUsername := chi.URLParam(r, "username")
	if targetUsername == "" || targetUsername == currentUsername {
		http.Error(w, "Invalid target", http.StatusBadRequest)
		return
	}

	receiver, err := repositories.FindUserByUsername(targetUsername)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	channel, err := repositories.FindChanBetweenTwoUsers(sender.ID, receiver.ID)
	if err != nil || channel == nil {
		channel, err = repositories.CreateChanBetweenTwoUsers(sender.ID, receiver.ID)
		if err != nil {
			http.Error(w, "Failed to create channel", http.StatusInternalServerError)
			return
		}
	}

	err = repositories.InsertChannelMsg(int(sender.ID), int(channel.ID), content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/channel/"+strconv.Itoa(int(channel.ID)))
	w.WriteHeader(http.StatusOK)
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
