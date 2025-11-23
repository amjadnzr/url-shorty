package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/amjadnzr/url-shortly/database"
	"github.com/amjadnzr/url-shortly/helpers"
	"github.com/amjadnzr/url-shortly/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	registerTemplate = "register"
	loginTemplate    = "login"
)

type Hanlder struct {
	db          *database.Database
	tokenHelper *helpers.TokenHelper
	templates   map[string]*template.Template
}

func NewHandler(db *database.Database, tokenHelper *helpers.TokenHelper) *Hanlder {
	templates := map[string]*template.Template{
		registerTemplate: template.Must(template.ParseFiles("templates/register.html")),
		loginTemplate:    template.Must(template.ParseFiles("templates/login.html")),
	}

	return &Hanlder{
		db:          db,
		tokenHelper: tokenHelper,
		templates:   templates,
	}
}

func (h *Hanlder) renderTemplate(w http.ResponseWriter, name string, data any) {
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "template rendering failed", http.StatusInternalServerError)
	}
}

func (h *Hanlder) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.renderTemplate(w, registerTemplate, nil)
}

func (h *Hanlder) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.renderTemplate(w, loginTemplate, nil)
}

// API
func (h *Hanlder) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := new(models.User)
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		fmt.Println(err)
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(passwordHash)

	id, err := h.db.CreateNewUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := h.tokenHelper.GenerateJWTToken(id)
	if err != nil {
		http.Error(w, "internal server error, try again later", http.StatusInternalServerError)
	}

	user.Token = token
	user.Id = id
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
	return
}

func (h *Hanlder) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Println("user   ", user)
	existingUser, err := h.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		fmt.Println("existing user ", existingUser, err)
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}
	fmt.Println("existing user", existingUser)
	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	token, err := h.tokenHelper.GenerateJWTToken(existingUser.Id)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "internal server error, try again later", http.StatusInternalServerError)
		return
	}
	existingUser.Token = token
	existingUser.Password = ""
	existingUser.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("HX-Redirect", "/register")
	_ = json.NewEncoder(w).Encode(existingUser)
	return
}
