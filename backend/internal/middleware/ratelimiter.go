package middleware

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	Enabled     bool
	MaxRequests int
	Window      time.Duration
}

type SlidingWindowRateLimiter struct {
	sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewSlidingWindowRateLimiter(limit int, window time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (rl *SlidingWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	rl.clients[ip] = append(rl.clients[ip], time.Now())

	// need to remove timestamps outside of window
	for i, t := range rl.clients[ip] {
		if time.Since(t) > rl.window {
			// ... allows Go to append each element in the provided slice one at a time
			rl.clients[ip] = append(rl.clients[ip][:i], rl.clients[ip][i+1:]...)
		}
	}

	// check if there are too many requests within the current window
	if len(rl.clients[ip]) > rl.limit {
		// this is the time when the oldest request in the window will expire
		retryAfter := rl.window - time.Since(rl.clients[ip][0])
		return false, retryAfter
	}

	return true, 0
}
