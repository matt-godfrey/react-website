package middleware

import (
	"fmt"
	"net/http"

	"github.com/matt-godfrey/react-website/internal/json"
)

func InternalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// w.WriteHeader(http.StatusInternalServerError)
	// _, _ = w.Write([]byte(err.Error()))
	json.Write(w, http.StatusInternalServerError, map[string]interface{}{
		"error": err.Error(),
	})
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter float64) {
	// w.WriteHeader(http.StatusTooManyRequests)
	// _, _ = w.Write([]byte("rate limit exceeded"))
	json.Write(w, http.StatusTooManyRequests, map[string]any{
		"error":       "rate limit exceeded",
		"retry_after": fmt.Sprintf("%.2f seconds", retryAfter),
	})
}
