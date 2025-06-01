package routes

import (
	"auth/components"
	"auth/handlers"
	"auth/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Setup(r chi.Router) {
	r.Use(middleware.SecurityHeaders)

	r.Get("/", handlers.Home)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/dashboard", handlers.HandleDashboard)
		r.Get("/getpgp", handlers.ServeGenPGP)
		r.Post("/getpgp", handlers.HandleGenPGP)

		// channels
		r.Get("/channels", handlers.GetOpenChannels)
		r.Get("/channels/explore", handlers.ServeChannelExplore)

		r.Get("/channel/{channel_id}/messages", handlers.HandleGetMessages)
		r.Post("/channel/{channel_id}/send", handlers.HandleSendMessage)

		r.Get("/channel/{channel_id}", handlers.ServeCommunication)

		// support route when two users don't have channels
		r.Get("/chat/{username}", handlers.ServeStartCommunication)
		r.Post("/chat/init/{username}", handlers.HandleChatInitSend)

		// utils
		r.Get("/online-users", handlers.HandleOnlineUsers)
		r.Get("/logout", handlers.HandleLogout)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthorizedMiddleware)

		r.Get("/login", handlers.ServeLogin)
		r.Get("/register", handlers.ServeRegister)
		r.Post("/login", handlers.HandleLogin)
		r.Post("/register", handlers.HandleRegister)
		r.Get("/register-success", handlers.RegisterSuccess)
	})

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		components.NotFoundPage().Render(r.Context(), w)
	}))
}
