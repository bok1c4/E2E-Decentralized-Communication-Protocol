package session

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var (
	store        *sessions.CookieStore
	session_name string
)

func InitSession() {
	session_name = os.Getenv("SESSION_NAME")
	if session_name == "" {
		session_name = "auth-session"
	}

	key := os.Getenv("SESSION_KEY")
	if key == "" {
		key = "some-key"
	}

	store = sessions.NewCookieStore([]byte(key))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,      // Prevent JS access
		Secure:   false,     // Change to true if using HTTPS
	}
}

// GetSession retrieves the session
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, session_name)
}

// SetSessionValue sets a value in the session
func SetSessionValue(w http.ResponseWriter, r *http.Request, key string, value any) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}

	session.Values[key] = value
	return session.Save(r, w)
}

// GetSessionValue retrieves a value from the session
func GetSessionValue(r *http.Request, key string) (any, error) {
	session, err := GetSession(r)
	if err != nil {
		return nil, err
	}
	return session.Values[key], nil
}

// DestroySession deletes the session
func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
