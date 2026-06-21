package middleware

import (
	"net"
	"net/http"
)

func RateLimiterMiddleware(rl RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				InternalErrorResponse(w, r, err)
				return
			}

			allow, retryAfter := rl.Allow(ip)
			if !allow {
				RateLimitExceededResponse(w, r, retryAfter)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ChiRateLimiterMiddleware(rl *ChiRateLimiter) func(http.Handler) http.Handler {
	return rl.rl.Handler
}
