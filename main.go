package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amjadnzr/url-shortly/database"
	"github.com/amjadnzr/url-shortly/handlers"
	"github.com/amjadnzr/url-shortly/helpers"
	"github.com/joho/godotenv"
)

// Question to be answered
// Why use mux instead of using http.NewServer and hanlders
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")

	db, err := database.InitDatabase(dbPath)
	if err != nil {
		log.Fatal("Database could not connect")
	}
	defer db.Close()

	tokenHelper := helpers.NewTokenHelper(jwtSecret)

	mux := http.NewServeMux()
	handler := handlers.NewHandler(db, tokenHelper)
	RegisterMux(mux, handler)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to run")
	}
}

func RegisterMux(m *http.ServeMux, h *handlers.Hanlder) {
	m.HandleFunc("/auth/register", h.RegisterUser)
	m.HandleFunc("/auth/login", h.LoginUser)
}
