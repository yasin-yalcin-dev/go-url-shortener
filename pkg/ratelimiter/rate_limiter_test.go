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
	"net/http/httptest"
	"testing"
)

// TestRateLimiter tests the basic functionality of rate limiting
func TestRateLimiter(t *testing.T) {
	// Create a new rate limiter
	limiter := NewRateLimiter(10, 20)

	testCases := []struct {
		name          string
		requestCount  int
		expectedAllow bool
	}{
		{
			name:          "Within Limit",
			requestCount:  5,
			expectedAllow: true,
		},
		{
			name:          "Exceeding Limit",
			requestCount:  25,
			expectedAllow: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip := "192.168.1.1"

			// Send multiple requests
			var lastAllowed bool
			for i := 0; i < tc.requestCount; i++ {
				lastAllowed = limiter.Allow(ip)
			}

			if lastAllowed != tc.expectedAllow {
				t.Errorf("Expected allow: %v, got %v", tc.expectedAllow, lastAllowed)
			}
		})
	}
}

// TestRateLimiterMiddleware tests the middleware functionality
func TestRateLimiterMiddleware(t *testing.T) {
	// Create a new rate limiter with very low limit
	limiter := NewRateLimiter(1, 1)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with rate limiter middleware
	limitedHandler := limiter.ChiMiddleware(testHandler)

	testCases := []struct {
		name               string
		requestCount       int
		expectedStatusCode int
	}{
		{
			name:               "First Request Allowed",
			requestCount:       1,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Subsequent Request Blocked",
			requestCount:       2,
			expectedStatusCode: http.StatusTooManyRequests,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var lastStatusCode int
			for i := 0; i < tc.requestCount; i++ {
				// Create a mock request
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.RemoteAddr = "192.168.1.2"

				// Create a response recorder
				w := httptest.NewRecorder()

				// Call the limited handler
				limitedHandler.ServeHTTP(w, req)

				lastStatusCode = w.Code
			}

			if lastStatusCode != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, lastStatusCode)
			}
		})
	}
}

// Benchmark rate limiter performance
func BenchmarkRateLimiter(b *testing.B) {
	limiter := NewRateLimiter(100, 200)
	ip := "192.168.1.3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow(ip)
	}
}
