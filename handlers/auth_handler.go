package handlers

import (
	"auth/components"
	"auth/services"
	"auth/session"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)
	page := components.Base("/register", isAuthenticated, components.RegisterForm())
	templ.Handler(page).ServeHTTP(w, r)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	pw := r.FormValue("pw")
	confirm_pw := r.FormValue("confirm_pw")

	if len(pw) <= 4 {
		http.Error(w, "Password must be at least 5 characters", http.StatusBadRequest)
		return
	}

	if pw != confirm_pw {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	err := services.RegisterUser(username, pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/register-success", http.StatusSeeOther)
}

func RegisterSuccess(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)
	page := components.Base("/register-success", isAuthenticated, components.RegistrationSuccess())
	if isAuthenticated {
		templ.Handler(page).ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := session.IsAuthenticated(r)

	page := components.Base("/login", isAuthenticated, components.LoginForm())
	templ.Handler(page).ServeHTTP(w, r)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	pw := r.FormValue("pw")

	err := services.LoginUser(username, pw)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, services.ErrWrongPassword) {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to log in: %v", err), http.StatusInternalServerError)
		return
	}

	err = session.SetSessionValue(w, r, "username", username)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/dashboard")
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := session.DestroySession(w, r)
	if err != nil {
		log.Printf("Failed to destroy session: %v", err)
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	err = session.SetSessionValue(w, r, "flash", "Successfully logged out.")
	if err != nil {
		log.Printf("Failed to set flash message: %v", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
