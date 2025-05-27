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
		r.Get("/chat", handlers.ServeIRC)
		r.Get("/chat/{username}", handlers.ServeDM)

		// TODO:
		// route for private messages
		// /chat/users_username, we need to protect this route
		// with validating if the user is in the database
		// (but we need safe operations attackers might try to
		// inject sql query)

		r.Get("/chat/messages", handlers.HandleGetMessages)
		r.Get("/online-users", handlers.HandleOnlineUsers)
		r.Post("/chat/send", handlers.HandleSendMessage)
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
