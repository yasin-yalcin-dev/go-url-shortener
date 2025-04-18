/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/yasin-yalcin-dev/go-url-shortener/internal/config"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/handler"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/service"
	customerrors "github.com/yasin-yalcin-dev/go-url-shortener/pkg/errors"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/logger"
)

// Mock URLStore for testing
type mockURLStore struct {
	urls map[string]string
}

func (m *mockURLStore) SaveShortenedURLWithTTL(ctx context.Context, shortID, originalURL string, ttl time.Duration) error {
	m.urls[shortID] = originalURL
	return nil
}

func (m *mockURLStore) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	url, exists := m.urls[shortID]
	if !exists {
		return "", customerrors.New(http.StatusNotFound, "URL not found")
	}
	return url, nil
}

func setupTestServer() (*handler.ShortenHandler, *chi.Mux) {
	// Create mock configuration
	cfg := &config.Config{
		BaseURL:       "http://localhost:8080",
		DefaultURLTTL: 24 * time.Hour,
	}

	// Create mock URL store
	mockStore := &mockURLStore{
		urls: make(map[string]string),
	}

	// Create mock logger
	mockLogger, err := logger.NewTestLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Create service with mock store
	urlService := service.NewURLShorteningService(cfg, mockStore)

	// Create handler
	handler := &handler.ShortenHandler{
		Service: urlService,
		Logger:  mockLogger,
	}

	// Create router
	r := chi.NewRouter()
	r.Post("/shorten", handler.ShortenURL)
	r.Get("/{shortened}", handler.Redirect)

	return handler, r
}

func TestFullURLShorteningWorkflow(t *testing.T) {
	// Setup test server
	_, router := setupTestServer()

	// Test URL to shorten
	originalURL := "https://example.com"

	// Prepare shorten request payload
	shortenPayload := map[string]string{
		"original": originalURL,
	}
	payloadBytes, _ := json.Marshal(shortenPayload)

	// Create request to shorten URL
	shortenReq := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(payloadBytes))
	shortenReq.Header.Set("Content-Type", "application/json")
	shortenW := httptest.NewRecorder()

	// Perform shorten request
	router.ServeHTTP(shortenW, shortenReq)

	// Check shorten response status
	assert.Equal(t, http.StatusOK, shortenW.Code)

	// Parse shortened URL from response
	var response map[string]string
	err := json.Unmarshal(shortenW.Body.Bytes(), &response)
	assert.NoError(t, err)

	shortenedURL := response["shortened"]
	assert.NotEmpty(t, shortenedURL)

	// Extract short ID
	shortID := shortenedURL[len(shortenedURL)-8:]

	// Create redirect request
	redirectReq := httptest.NewRequest("GET", "/"+shortID, nil)
	redirectW := httptest.NewRecorder()

	// Perform redirect
	router.ServeHTTP(redirectW, redirectReq)

	// Check redirect response
	assert.Equal(t, http.StatusFound, redirectW.Code)
	assert.Equal(t, originalURL, redirectW.Header().Get("Location"))
}
