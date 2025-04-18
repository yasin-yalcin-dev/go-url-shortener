/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLStore interface {

	// SaveShortenedURLWithTTL stores a shortened URL with a time-to-live (TTL) in the database
	SaveShortenedURLWithTTL(ctx context.Context, shortID, originalURL string, ttl time.Duration) error
	// GetOriginalURL retrieves the original URL from the database
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

// RedisStore struct implements the URLStore interface for Redis.
type RedisStore struct {
	Client *redis.Client
}

// NewRedisStore creates a new RedisStore instance
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{Client: client}
}

// GetOriginalURL retrieves the original URL from Redis
func (r *RedisStore) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	originalURL, err := r.Client.Get(ctx, shortID).Result()
	if err != nil {
		return "", fmt.Errorf("could not get original URL: %v", err)
	}
	return originalURL, nil
}

func (r *RedisStore) SaveShortenedURLWithTTL(
	ctx context.Context,
	shortID,
	originalURL string,
	ttl time.Duration,
) error {
	err := r.Client.Set(ctx, shortID, originalURL, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save URL with TTL: %v", err)
	}
	return nil
}
