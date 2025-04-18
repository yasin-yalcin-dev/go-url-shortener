/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package service

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yasin-yalcin-dev/go-url-shortener/internal/config"
)

// Mock Redis Store
type mockRedisStore struct {
	urls map[string]string
}

func (m *mockRedisStore) SaveShortenedURLWithTTL(ctx context.Context, shortID, originalURL string, ttl time.Duration) error {
	m.urls[shortID] = originalURL
	return nil
}

func (m *mockRedisStore) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	url, exists := m.urls[shortID]
	if !exists {
		return "", fmt.Errorf("URL not found")
	}
	return url, nil
}

func TestShortenURL(t *testing.T) {
	testCases := []struct {
		name          string
		originalURL   string
		expectedError bool
		ttlOption     time.Duration
	}{
		{
			name:          "Valid URL",
			originalURL:   "https://example.com",
			expectedError: false,
		},
		{
			name:          "Invalid URL",
			originalURL:   "not-a-url",
			expectedError: true,
		},
		{
			name:          "Empty URL",
			originalURL:   "",
			expectedError: true,
		},
		{
			name:          "URL with Custom TTL",
			originalURL:   "https://example.com",
			expectedError: false,
			ttlOption:     1 * time.Hour,
		},
	}

	// Create configuration and mock store
	cfg := &config.Config{
		BaseURL:       "http://short.url",
		DefaultURLTTL: 24 * time.Hour,
	}
	mockStore := &mockRedisStore{urls: make(map[string]string)}

	// Create the URL shortening service
	service := NewURLShorteningService(cfg, mockStore)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// URL shortening
			var shortenedURL string
			var err error

			if tc.ttlOption > 0 {
				shortenedURL, err = service.ShortenURL(context.Background(), tc.originalURL, WithTTL(tc.ttlOption))
			} else {
				shortenedURL, err = service.ShortenURL(context.Background(), tc.originalURL)
			}

			// Error handling
			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// URL format control
				if !strings.HasPrefix(shortenedURL, cfg.BaseURL) {
					t.Errorf("Shortened URL does not start with base URL")
				}
			}
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	// create configuration and mock store
	cfg := &config.Config{
		BaseURL:       "http://short.url",
		DefaultURLTTL: 24 * time.Hour,
	}
	mockStore := &mockRedisStore{
		urls: map[string]string{
			"existing-id": "https://example.com",
		},
	}

	// Screate the URL shortening service
	service := NewURLShorteningService(cfg, mockStore)

	testCases := []struct {
		name          string
		shortID       string
		expectedError bool
	}{
		{
			name:          "Existing URL",
			shortID:       "existing-id",
			expectedError: false,
		},
		{
			name:          "Non-Existing URL",
			shortID:       "non-existing-id",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalURL, err := service.GetOriginalURL(context.Background(), tc.shortID)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// URL control
				if originalURL != "https://example.com" {
					t.Errorf("Unexpected original URL")
				}
			}
		})
	}
}

// Performans test for URL shortening
func BenchmarkShortenURL(b *testing.B) {
	cfg := &config.Config{
		BaseURL:       "http://short.url",
		DefaultURLTTL: 24 * time.Hour,
	}
	mockStore := &mockRedisStore{urls: make(map[string]string)}
	service := NewURLShorteningService(cfg, mockStore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ShortenURL(context.Background(), "https://example.com")
	}
}
