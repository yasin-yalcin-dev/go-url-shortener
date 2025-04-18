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
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/config"
)

// RedisClient represents an enhanced Redis client
type RedisClient struct {
	client *redis.Client
	config *config.RedisConfig
}

// Connect establishes a connection to Redis and configures the client
func Connect(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.RedisConfig.Address,
		Password:        cfg.RedisConfig.Password,
		DB:              cfg.RedisConfig.DB,
		PoolSize:        cfg.RedisConfig.PoolSize,
		DialTimeout:     cfg.RedisConfig.DialTimeout,
		ReadTimeout:     cfg.RedisConfig.ReadTimeout,
		WriteTimeout:    cfg.RedisConfig.WriteTimeout,
		PoolTimeout:     cfg.RedisConfig.PoolTimeout,
		MaxRetries:      cfg.RedisConfig.MaxRetries,
		MinRetryBackoff: cfg.RedisConfig.MinRetryBackoff,
		MaxRetryBackoff: cfg.RedisConfig.MaxRetryBackoff,
	})

	// Check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection error: %v", err)
	}

	log.Println("Redis connection successful!")

	return &RedisClient{
		client: client,
		config: cfg.RedisConfig,
	}, nil
}

// Client returns the raw Redis client
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// Close safely closes the connection
func (r *RedisClient) Close() error {
	if err := r.client.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
		return err
	}
	return nil
}

// IsHealthy checks if the Redis connection is healthy
func (r *RedisClient) IsHealthy(ctx context.Context) bool {
	_, err := r.client.Ping(ctx).Result()
	return err == nil
}

// GetConfig returns the current Redis configuration
func (r *RedisClient) GetConfig() *config.RedisConfig {
	return r.config
}

// ReconnectWithBackoff attempts to reconnect with exponential backoff
func (r *RedisClient) ReconnectWithBackoff(cfg *config.Config) error {
	backoff := time.Second
	maxBackoff := 1 * time.Minute

	for attempt := 1; attempt <= 5; attempt++ {
		newClient, err := Connect(cfg)
		if err == nil {
			r.client = newClient.client
			r.config = newClient.config
			return nil
		}

		log.Printf("Reconnection attempt %d failed: %v", attempt, err)

		time.Sleep(backoff)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	return fmt.Errorf("5 reconnection attempts failed")
}
