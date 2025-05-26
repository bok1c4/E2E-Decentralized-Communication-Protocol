package middleware

import (
	"auth/session"
	"net/http"
)

// prevent access from non-authorized users (protected routes)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, err := session.GetSessionValue(r, "username")
		if err != nil || username == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Redirects to home if user is already authenticated
func AuthorizedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated := session.IsAuthenticated(r)

		if isAuthenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	// download htmx and serve it by myself
	// Strict CSP — HTMX must be served locally

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline';")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains") // Only if HTTPS

		next.ServeHTTP(w, r)
	})
}
