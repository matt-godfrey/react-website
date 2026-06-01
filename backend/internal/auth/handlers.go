package auth

import (
	stdjson "encoding/json"
	"log"
	"net/http"

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

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	err := h.service.RegisterUser(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := stdjson.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("%s %s %s", creds.Username, creds.Email, creds.Password)
	json.Write(w, http.StatusOK, creds)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

}
