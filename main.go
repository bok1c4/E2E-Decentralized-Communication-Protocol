package main

import (
	"auth/config"
	"auth/db"
	"auth/routes"
	"auth/session"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db.InitDB(db.Config(cfg.Database))
	session.InitSession()

	r := chi.NewRouter()
	routes.Setup(r)

	log.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
