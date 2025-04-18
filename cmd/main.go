/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/yasin-yalcin-dev/go-url-shortener/internal/config"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/handler"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/redis"
	"github.com/yasin-yalcin-dev/go-url-shortener/internal/service"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/analytics"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/logger"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/ratelimiter"

	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Create logger
	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	// Create rate limiter
	// 10 requests per second, burst of 20 requests
	rateLimiter := ratelimiter.NewRateLimiter(10, 20)

	// Clean up old entries every hour
	rateLimiter.Clean(1 * time.Hour)

	// Connect to Redis
	redisClient, err := redis.Connect(cfg)
	if err != nil {
		appLogger.Error("Redis connection failed",
			zap.Error(err),
			zap.String("address", cfg.RedisConfig.Address),
		)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			appLogger.Error("Failed to close Redis connection",
				zap.Error(err),
				zap.String("address", cfg.RedisConfig.Address),
			)
		}
	}()

	// Initialize Redis store
	redisStore := redis.NewRedisStore(redisClient.Client())

	// Initialize service
	urlService := service.NewURLShorteningService(cfg, redisStore)

	// Analytics store
	analyticsStore := analytics.NewAnalyticsStore(redisClient.Client())

	// Initialize handler
	shortenHandler := &handler.ShortenHandler{
		Service:   urlService,
		Logger:    appLogger,
		Analytics: analyticsStore,
	}
	// Create a new router
	r := chi.NewRouter()

	// Apply rate limiting middleware
	r.Use(rateLimiter.ChiMiddleware)

	// Routes
	r.Post("/shorten", shortenHandler.ShortenURL)  // URL shortening endpoint
	r.Get("/{shortened}", shortenHandler.Redirect) // URL redirect endpoint
	r.Get("/{shortened}/analytics", shortenHandler.GetURLAnalytics)

	// Start the server
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		appLogger.Error("Server startup failed", zap.Error(err))
		log.Fatalf("Server startup failed: %v", err)
	}
}
