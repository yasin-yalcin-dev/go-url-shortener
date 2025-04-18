/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// RedisConfig represents the configuration for Redis connection
type RedisConfig struct {
	Address         string        // Redis server address (host:port)
	Password        string        // Redis authentication password
	DB              int           // Database number to use
	PoolSize        int           // Connection pool size
	DialTimeout     time.Duration // Timeout for establishing connection
	ReadTimeout     time.Duration // Timeout for read operations
	WriteTimeout    time.Duration // Timeout for write operations
	PoolTimeout     time.Duration // Timeout for connection pool
	MaxRetries      int           // Maximum number of retries
	MinRetryBackoff time.Duration // Minimum backoff time between retries
	MaxRetryBackoff time.Duration // Maximum backoff time between retries
}

// Config holds the overall application configuration
type Config struct {
	RedisConfig   *RedisConfig
	ServerPort    string
	BaseURL       string
	LogLevel      string
	DefaultURLTTL time.Duration
}

// Load Loads the .env file and environment variables
func Load() (*Config, error) {
	// Load .env file, otherwise use environment variables
	godotenv.Load()

	cfg := &Config{
		RedisConfig:   defaultRedisConfig(),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		DefaultURLTTL: getDurationEnv("DEFAULT_URL_TTL", 24*time.Hour),
	}
	// verify configuration
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// defaultRedisConfig creates default Redis configuration values
func defaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Address:         getEnv("REDIS_ADDR", "localhost:6379"),
		Password:        getEnv("REDIS_PASSWORD", ""),
		DB:              getEnvAsInt("REDIS_DB", 0),
		PoolSize:        getEnvAsInt("REDIS_POOL_SIZE", 10),
		DialTimeout:     getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:     getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout:    getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		PoolTimeout:     getEnvAsDuration("REDIS_POOL_TIMEOUT", 4*time.Second),
		MaxRetries:      getEnvAsInt("REDIS_MAX_RETRIES", 3),
		MinRetryBackoff: getEnvAsDuration("REDIS_MIN_RETRY_BACKOFF", 300*time.Millisecond),
		MaxRetryBackoff: getEnvAsDuration("REDIS_MAX_RETRY_BACKOFF", 2*time.Second),
	}
}

// validate checks if the configuration is valid
func validate(cfg *Config) error {
	// Validate Redis address
	if cfg.RedisConfig.Address == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}

	// Validate server port
	if cfg.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	// Validate base URL
	if cfg.BaseURL == "" {
		return fmt.Errorf("BASE_URL is required")
	}
	return nil
}

// getEnv retrieves an environment variable, returns default if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt converts environment variable to integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getEnvAsDuration converts environment variable to time.Duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return duration
}
