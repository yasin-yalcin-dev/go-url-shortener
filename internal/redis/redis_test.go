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
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// We will use miniredis to test without real Redis connection
func setupMockRedis() (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestRedisStore_SaveAndGetShortenedURL(t *testing.T) {
	// setup mock Redis
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	// create RedisStore
	store := NewRedisStore(client)

	// Test cases
	testCases := []struct {
		name        string
		shortID     string
		originalURL string
		ttl         time.Duration
	}{
		{
			name:        "Save and Retrieve URL",
			shortID:     "test123",
			originalURL: "https://example.com",
			ttl:         24 * time.Hour,
		},
		{
			name:        "Save and Retrieve URL with No TTL",
			shortID:     "test456",
			originalURL: "https://another-example.com",
			ttl:         0,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// save the shortened URL
			err := store.SaveShortenedURLWithTTL(ctx, tc.shortID, tc.originalURL, tc.ttl)
			if err != nil {
				t.Fatalf("SaveShortenedURLWithTTL failed: %v", err)
			}

			// retrieve the original URL
			retrievedURL, err := store.GetOriginalURL(ctx, tc.shortID)
			if err != nil {
				t.Fatalf("GetOriginalURL failed: %v", err)
			}

			// verification
			if retrievedURL != tc.originalURL {
				t.Errorf("Expected URL %s, got %s", tc.originalURL, retrievedURL)
			}
		})
	}
}

func TestRedisStore_URLExpiration(t *testing.T) {
	// setup mock Redis
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	// create RedisStore
	store := NewRedisStore(client)

	ctx := context.Background()
	shortID := "expire-test"
	originalURL := "https://example.com"
	shortTTL := 1 * time.Second

	// save the shortened URL with TTL
	err := store.SaveShortenedURLWithTTL(ctx, shortID, originalURL, shortTTL)
	if err != nil {
		t.Fatalf("SaveShortenedURLWithTTL failed: %v", err)
	}

	// retrieve the original URL immediately
	_, err = store.GetOriginalURL(ctx, shortID)
	if err != nil {
		t.Fatalf("GetOriginalURL failed immediately: %v", err)
	}

	// Fast forward time instead of sleeping
	mr.FastForward(2 * time.Second)

	// try to retrieve the original URL after expiration
	_, err = store.GetOriginalURL(ctx, shortID)
	if err == nil {
		t.Errorf("Expected URL to expire, but it still exists ")
	}
}

func TestRedisStore_NonExistentURL(t *testing.T) {
	// setup mock Redis
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	// create RedisStore
	store := NewRedisStore(client)

	ctx := context.Background()
	nonExistentID := "non-existent-id"

	// try to retrieve a non-existent URL
	_, err := store.GetOriginalURL(ctx, nonExistentID)
	if err == nil {
		t.Errorf("Expected error when retrieving non-existent URL")
	}
}

// performance test for URL shortening
func BenchmarkRedisStore_SaveAndGet(b *testing.B) {
	mr, client := setupMockRedis()
	defer mr.Close()
	defer client.Close()

	store := NewRedisStore(client)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shortID := fmt.Sprintf("bench-%d", i)
		originalURL := fmt.Sprintf("https://example-%d.com", i)

		// save the shortened URL
		err := store.SaveShortenedURLWithTTL(ctx, shortID, originalURL, 0)
		if err != nil {
			b.Fatalf("SaveShortenedURLWithTTL failed: %v", err)
		}

		// get the original URL
		_, err = store.GetOriginalURL(ctx, shortID)
		if err != nil {
			b.Fatalf("GetOriginalURL failed: %v", err)
		}
	}
}
