package main

import (
	"log"
	"net/http"

	"github.com/amjadnzr/url-shortly/database"
	"github.com/amjadnzr/url-shortly/handlers"
)

// Question to be answered
// Why use mux instead of using http.NewServer and hanlders
func main() {
	mux := http.NewServeMux()

	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal("Database could not connect")
	}
	defer db.Close()

	handler := handlers.NewHandler(db)
	RegisterMux(mux, handler)

	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal("Server failed to run")
	}
}

func RegisterMux(m *http.ServeMux, h *handlers.Hanlder) {
	m.HandleFunc("/auth/register", h.RegisterUser)
	m.HandleFunc("/auth/login", h.LoginUser)
}
