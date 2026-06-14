package quotes

import (
	"net/http"

	"github.com/matt-godfrey/react-website/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.service.GetRandomQuote(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, quote)
}
