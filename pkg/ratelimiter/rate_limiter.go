/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package ratelimiter

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter controls request rates for different clients
type RateLimiter struct {
	visitors map[string]*visitorState
	mutex    sync.Mutex
	limit    rate.Limit
	burst    int
}

// visitorState IP bazlı rate limiter durumu
type visitorState struct {
	limiter    *rate.Limiter
	lastActive time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitorState),
		limit:    rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// Allow checks if a request from a specific IP is allowed
func (r *RateLimiter) Allow(ip string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Create a new limiter for the IP if it doesn't exist
	if _, exists := r.visitors[ip]; !exists {
		r.visitors[ip] = &visitorState{
			limiter:    rate.NewLimiter(r.limit, r.burst),
			lastActive: time.Now(),
		}
	}

	// Zaman güncellemesi
	visitor := r.visitors[ip]
	visitor.lastActive = time.Now()

	// Check if the request is allowed
	return visitor.limiter.Allow()
}

// Middleware provides HTTP middleware for rate limiting
func (r *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get client IP address
		ip := getIP(req)

		// Check rate limit
		if !r.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, req)
	}
}

// getIP extracts the real IP address from the request
func getIP(req *http.Request) string {
	// Check for forwarded IP (useful behind proxies)
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return ip
}

// Clean periodically removes old entries to prevent memory leaks
func (r *RateLimiter) Clean(duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			r.mutex.Lock()
			for ip, visitor := range r.visitors {
				// Remove entries that haven't been used recently
				if time.Since(visitor.lastActive) > duration {
					delete(r.visitors, ip)
				}
			}
			r.mutex.Unlock()
		}
	}()
}

// ChiMiddleware Chi router için uyumlu middleware
func (r *RateLimiter) ChiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get client IP address
		ip := getIP(req)

		// Check rate limit
		if !r.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, req)
	})
}
