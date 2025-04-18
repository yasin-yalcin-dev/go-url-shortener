/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/model"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/service"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/analytics"
	customerrors "github.com/yasin-yalcin-dev/go-url-shortener/pkg/errors"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/logger"
)

// mockURLService simulates the URL shortening service for testing
type mockURLService struct {
	shortenFunc     func(ctx context.Context, url string, options ...service.URLShortenOption) (string, error)
	getOriginalFunc func(ctx context.Context, shortID string) (string, error)
}

// ShortenURL implements the URL shortening method for the mock service
func (m *mockURLService) ShortenURL(ctx context.Context, url string, options ...service.URLShortenOption) (string, error) {
	return m.shortenFunc(ctx, url, options...)
}

// mockAnalyticsStore simulates the analytics store for testing
type mockAnalyticsStore struct {
	recordFunc func(ctx context.Context, shortID, ipAddress string) error
	getFunc    func(ctx context.Context, shortID string) (*analytics.URLAnalytics, error)
}

// RecordURLAccess implements the URL access recording method for the mock analytics store
func (m *mockAnalyticsStore) RecordURLAccess(ctx context.Context, shortID, ipAddress string) error {
	if m.recordFunc != nil {
		return m.recordFunc(ctx, shortID, ipAddress)
	}
	return nil
}

// GetURLAnalytics implements the URL analytics retrieval method
func (m *mockAnalyticsStore) GetURLAnalytics(ctx context.Context, shortID string) (*analytics.URLAnalytics, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, shortID)
	}
	return nil, nil
}

// GetOriginalURL implements the original URL retrieval method for the mock service
func (m *mockURLService) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	return m.getOriginalFunc(ctx, shortID)
}

// setUp prepares the test environment
func setUp(t *testing.T) {
	// Ensure logs directory exists
	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}
}

func TestShortenHandler_ShortenURL(t *testing.T) {
	// Prepare test environment
	setUp(t)

	// Test cases for URL shortening
	testCases := []struct {
		name               string
		requestBody        model.URL
		mockShortenFunc    func(ctx context.Context, url string, options ...service.URLShortenOption) (string, error)
		expectedStatusCode int
	}{
		{
			name: "Successful URL Shortening",
			requestBody: model.URL{
				Original: "https://example.com",
			},
			mockShortenFunc: func(ctx context.Context, url string, options ...service.URLShortenOption) (string, error) {
				return "http://short.url/abc123", nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid URL",
			requestBody: model.URL{
				Original: "invalid-url",
			},
			mockShortenFunc: func(ctx context.Context, url string, options ...service.URLShortenOption) (string, error) {
				return "", customerrors.New(
					http.StatusBadRequest,
					"Invalid URL",
					"URL format is incorrect",
				)
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := &mockURLService{
				shortenFunc: tc.mockShortenFunc,
			}

			// Create mock logger
			mockLogger, err := logger.New("info")
			if err != nil {
				t.Fatalf("Failed to create mock logger: %v", err)
			}

			// Create handler
			handler := &ShortenHandler{
				Service: mockService,
				Logger:  mockLogger,
			}

			// Prepare request body
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.ShortenURL(w, req)

			// Check response status
			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tc.expectedStatusCode, w.Code)
			}

			// Check response content
			if tc.expectedStatusCode == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if response["original"] != tc.requestBody.Original {
					t.Errorf("Expected original URL %s, got %s", tc.requestBody.Original, response["original"])
				}
				if response["shortened"] == "" {
					t.Errorf("Expected non-empty shortened URL")
				}
			}
		})
	}
}

func TestShortenHandler_Redirect(t *testing.T) {
	// Prepare test environment
	setUp(t)

	// Test cases for URL redirection
	testCases := []struct {
		name               string
		shortID            string
		mockGetOriginalURL func(ctx context.Context, shortID string) (string, error)
		mockRecordAccess   func(ctx context.Context, shortID, ipAddress string) error
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name:    "Successful Redirect",
			shortID: "abc123",
			mockGetOriginalURL: func(ctx context.Context, shortID string) (string, error) {
				return "https://example.com", nil
			},
			mockRecordAccess: func(ctx context.Context, shortID, ipAddress string) error {
				return nil
			},
			expectedStatusCode: http.StatusFound,
			expectedLocation:   "https://example.com",
		},
		{
			name:    "URL Not Found",
			shortID: "non-existent",
			mockGetOriginalURL: func(ctx context.Context, shortID string) (string, error) {
				return "", customerrors.New(
					http.StatusNotFound,
					"Short URL not found",
					"The requested short URL does not exist",
				)
			},
			mockRecordAccess: func(ctx context.Context, shortID, ipAddress string) error {
				return nil
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := &mockURLService{
				getOriginalFunc: tc.mockGetOriginalURL,
			}

			// Create mock analytics store
			mockAnalytics := &mockAnalyticsStore{
				recordFunc: tc.mockRecordAccess,
				getFunc: func(ctx context.Context, shortID string) (*analytics.URLAnalytics, error) {
					return &analytics.URLAnalytics{}, nil
				},
			}

			// Create mock logger
			mockLogger, err := logger.New("info")
			if err != nil {
				t.Fatalf("Failed to create mock logger: %v", err)
			}

			// Create handler
			handler := &ShortenHandler{
				Service:   mockService,
				Logger:    mockLogger,
				Analytics: mockAnalytics,
			}

			// Create Chi router context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shortened", tc.shortID)

			// Create request
			req, _ := http.NewRequest("GET", "/"+tc.shortID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Redirect(w, req)

			// Check response status
			if w.Code != tc.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tc.expectedStatusCode, w.Code)
			}

			// Check redirection for successful case
			if tc.expectedStatusCode == http.StatusFound {
				location := w.Header().Get("Location")
				if location != tc.expectedLocation {
					t.Errorf("Expected redirect to %s, got %s", tc.expectedLocation, location)
				}
			}
		})
	}
}
