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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yasin-yalcin-dev/go-url-shortener/internal/service"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/analytics"
	customerrors "github.com/yasin-yalcin-dev/go-url-shortener/pkg/errors"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/logger"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type ShortenHandler struct {
	Service   service.URLShorteningService
	Logger    *logger.Logger
	Analytics analytics.AnalyticsStoreInterface
}

// ShortenURL will create a shortened URL
func (h *ShortenHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var urlRequest struct {
		Original string        `json:"original"`
		TTL      time.Duration `json:"ttl,omitempty"`
	}

	// Log incoming request
	h.Logger.Info("Received URL shortening request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&urlRequest); err != nil {
		h.Logger.Error("Failed to decode request body",
			zap.Error(err),
			zap.String("body", fmt.Sprintf("%+v", r.Body)),
		)
		apiErr := customerrors.New(
			http.StatusBadRequest,
			"Invalid input",
			err.Error(),
		)
		apiErr.WriteResponse(w)
		return
	}

	// Shorten URL
	var shortenedURL string
	var err error
	if urlRequest.TTL > 0 {
		shortenedURL, err = h.Service.ShortenURL(
			r.Context(),
			urlRequest.Original,
			service.WithTTL(urlRequest.TTL),
		)
	} else {
		shortenedURL, err = h.Service.ShortenURL(r.Context(), urlRequest.Original)
	}
	if err != nil {
		h.Logger.Error("URL shortening failed",
			zap.Error(err),
			zap.String("originalURL", urlRequest.Original),
		)
		// Check if it's an APIError
		if apiErr, ok := err.(*customerrors.APIError); ok {
			apiErr.WriteResponse(w)
		} else {
			// Fallback to internal server error
			customerrors.ErrInternal.WriteResponse(w)
		}
		return
	}

	// Log successful shortening
	h.Logger.Info("URL successfully shortened",
		zap.String("originalURL", urlRequest.Original),
		zap.String("shortenedURL", shortenedURL),
	)

	response := map[string]string{
		"original":  urlRequest.Original,
		"shortened": shortenedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Redirect will handle redirection from short URL to the original URL
func (h *ShortenHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Get the short ID from the URL
	shortID := chi.URLParam(r, "shortened")

	// Log redirect attempt
	h.Logger.Info("Redirect attempt",
		zap.String("shortID", shortID),
	)
	// Fetch original URL from Redis
	originalURL, err := h.Service.GetOriginalURL(r.Context(), shortID)
	if err != nil {
		h.Logger.Error("URL redirect failed",
			zap.Error(err),
			zap.String("shortID", shortID),
		)

		// Handle URL not found
		apiErr := customerrors.New(
			http.StatusNotFound,
			"Short URL not found",
			"The requested short URL does not exist",
		)
		apiErr.WriteResponse(w)
		return
	}
	// Save analytics
	go func() {
		// Get client IP
		ip := getClientIP(r)

		// Save URL access analytics
		if err := h.Analytics.RecordURLAccess(context.Background(), shortID, ip); err != nil {
			h.Logger.Error("Failed to record URL access",
				zap.Error(err),
				zap.String("shortID", shortID),
			)
		}
	}()
	// Log successful redirect
	h.Logger.Info("Successful redirect",
		zap.String("shortID", shortID),
		zap.String("originalURL", originalURL),
	)
	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *ShortenHandler) GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortened")

	// get analytics from the store
	analytics, err := h.Analytics.GetURLAnalytics(r.Context(), shortID)
	if err != nil {
		h.Logger.Error("Failed to get URL analytics",
			zap.Error(err),
			zap.String("shortID", shortID),
		)
		customerrors.ErrInternal.WriteResponse(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// getClientIP get the client IP address from the request
func getClientIP(r *http.Request) string {
	// Use X-Forwarded-For or X-Real-IP if available
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
