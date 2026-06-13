package auth

import (
	stdjson "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/matt-godfrey/react-website/internal/json"
)

// type RegisterRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

//	type LoginRequest struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// Register handles the registration request
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := stdjson.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.service.RegisterUser(r.Context(), creds.Username, creds.Email, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("%s %s %s", creds.Username, creds.Email, creds.Password)
	json.Write(w, http.StatusOK, creds)
}

// Login handles the login request
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := stdjson.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	sessionId, err := h.service.LoginUser(r.Context(), creds.Email, creds.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Secure:   false, // change for prod
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	json.Write(w, http.StatusOK, sessionId)
}

// Logout handles the logout request
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session cookie", http.StatusUnauthorized)
		return
	}
	sessionId := cookie.Value

	err = h.service.LogoutUser(r.Context(), sessionId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	json.Write(w, http.StatusOK, sessionId)
}

// GetCurrentUser handles the get current user request for a session
// It returns the current user if the session is valid, or an error if not
func (h *handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session cookie", http.StatusUnauthorized)
		return
	}
	sessionId := cookie.Value

	// if the session is valid, return the current user
	user, err := h.service.Authenticate(r.Context(), sessionId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.Write(w, http.StatusOK, user)
}
