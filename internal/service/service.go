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
	"time"

	"github.com/yasin-yalcin-dev/go-url-shortener/internal/config"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/redis"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/shortener"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/validator"
)

// URLShorteningService defines methods for URL shortening
type URLShorteningService interface {
	ShortenURL(ctx context.Context, originalURL string, options ...URLShortenOption) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

// URLShorteningServiceImpl implements the URLShorteningService interface
type URLShorteningServiceImpl struct {
	cfg       *config.Config
	Store     redis.URLStore
	validator *validator.URLValidator
}

func NewURLShorteningService(cfg *config.Config, store redis.URLStore) *URLShorteningServiceImpl {
	return &URLShorteningServiceImpl{
		cfg:       cfg,
		Store:     store,
		validator: validator.NewURLValidator(),
	}
}

// Structures for URL shortening options
type URLShortenOption func(*urlShortenOptions)

type urlShortenOptions struct {
	ttl time.Duration
}

// Optional function to set expiration time withTTL
func WithTTL(duration time.Duration) URLShortenOption {
	return func(opts *urlShortenOptions) {
		opts.ttl = duration
	}
}

func (s *URLShorteningServiceImpl) ShortenURL(ctx context.Context, originalURL string, options ...URLShortenOption) (string, error) {
	// Validate URL
	if apiErr := s.validator.Validate(originalURL); apiErr != nil {
		return "", apiErr
	}

	opts := &urlShortenOptions{
		ttl: s.cfg.DefaultURLTTL, // Varsayılan süre
	}

	for _, opt := range options {
		opt(opts)
	}

	shortID, err := shortener.GenerateUnique(func(id string) bool {
		// Check if this ID exists in Redis
		_, err := s.Store.GetOriginalURL(ctx, id)
		return err == nil // If the error is nil, the ID already exists
	})
	if err != nil {
		return "", err
	}
	err = s.Store.SaveShortenedURLWithTTL(ctx, shortID, originalURL, opts.ttl)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortID), nil
}

func (s *URLShorteningServiceImpl) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	return s.Store.GetOriginalURL(ctx, shortID)
}
