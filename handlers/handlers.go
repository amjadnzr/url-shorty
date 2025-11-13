package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/amjadnzr/url-shortly/database"
	"github.com/amjadnzr/url-shortly/models"
	"golang.org/x/crypto/bcrypt"
)

type Hanlder struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Hanlder {
	return &Hanlder{
		db: db,
	}
}

func (h *Hanlder) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}
	user.PasswordHash = string(passwordHash)

	id, err := h.db.CreateNewUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	user.Id = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	return

}

func (h *Hanlder) LoginUser(w http.ResponseWriter, r *http.Request) {

}
