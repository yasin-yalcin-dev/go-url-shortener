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
	"os"
	"testing"
	"time"
)

// Helper function to temporarily set environment variables
func setEnv(key, value string) func() {
	oldValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, oldValue)
	}
}

func TestLoad(t *testing.T) {
	testCases := []struct {
		name           string
		envSetup       func() func()
		expectedConfig *Config
	}{
		{
			name: "Default Configuration",
			envSetup: func() func() {
				return func() {}
			},
			expectedConfig: &Config{
				RedisConfig: &RedisConfig{
					Address:         "localhost:6379",
					Password:        "",
					DB:              0,
					PoolSize:        10,
					DialTimeout:     5 * time.Second,
					ReadTimeout:     3 * time.Second,
					WriteTimeout:    3 * time.Second,
					PoolTimeout:     4 * time.Second,
					MaxRetries:      3,
					MinRetryBackoff: 300 * time.Millisecond,
					MaxRetryBackoff: 2 * time.Second,
				},
				ServerPort:    "8080",
				BaseURL:       "http://localhost:8080",
				LogLevel:      "info",
				DefaultURLTTL: 24 * time.Hour,
			},
		},
		{
			name: "Custom Configuration",
			envSetup: func() func() {
				resetFuncs := []func(){}
				resetFuncs = append(resetFuncs, setEnv("REDIS_ADDR", "custom-redis:6380"))
				resetFuncs = append(resetFuncs, setEnv("SERVER_PORT", "9090"))
				resetFuncs = append(resetFuncs, setEnv("BASE_URL", "http://custom.url"))
				resetFuncs = append(resetFuncs, setEnv("LOG_LEVEL", "debug"))
				resetFuncs = append(resetFuncs, setEnv("DEFAULT_URL_TTL", "48h"))

				return func() {
					for _, reset := range resetFuncs {
						reset()
					}
				}
			},
			expectedConfig: &Config{
				RedisConfig: &RedisConfig{
					Address:         "custom-redis:6380",
					Password:        "",
					DB:              0,
					PoolSize:        10,
					DialTimeout:     5 * time.Second,
					ReadTimeout:     3 * time.Second,
					WriteTimeout:    3 * time.Second,
					PoolTimeout:     4 * time.Second,
					MaxRetries:      3,
					MinRetryBackoff: 300 * time.Millisecond,
					MaxRetryBackoff: 2 * time.Second,
				},
				ServerPort:    "9090",
				BaseURL:       "http://custom.url",
				LogLevel:      "debug",
				DefaultURLTTL: 48 * time.Hour,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			resetEnv := tc.envSetup()
			defer resetEnv()

			// Load Configuration
			cfg, err := Load()
			if err != nil {
				t.Fatalf("Unexpected error loading config: %v", err)
			}

			// Redis configuration check
			if cfg.RedisConfig.Address != tc.expectedConfig.RedisConfig.Address {
				t.Errorf("Expected Redis Address %s, got %s",
					tc.expectedConfig.RedisConfig.Address,
					cfg.RedisConfig.Address)
			}

			// Check other configuration fields
			if cfg.ServerPort != tc.expectedConfig.ServerPort {
				t.Errorf("Expected ServerPort %s, got %s",
					tc.expectedConfig.ServerPort,
					cfg.ServerPort)
			}

			if cfg.BaseURL != tc.expectedConfig.BaseURL {
				t.Errorf("Expected BaseURL %s, got %s",
					tc.expectedConfig.BaseURL,
					cfg.BaseURL)
			}

			if cfg.LogLevel != tc.expectedConfig.LogLevel {
				t.Errorf("Expected LogLevel %s, got %s",
					tc.expectedConfig.LogLevel,
					cfg.LogLevel)
			}

			if cfg.DefaultURLTTL != tc.expectedConfig.DefaultURLTTL {
				t.Errorf("Expected DefaultURLTTL %v, got %v",
					tc.expectedConfig.DefaultURLTTL,
					cfg.DefaultURLTTL)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "Valid Configuration",
			config: &Config{
				RedisConfig: &RedisConfig{Address: "localhost:6379"},
				ServerPort:  "8080",
				BaseURL:     "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "Missing Redis Address",
			config: &Config{
				RedisConfig: &RedisConfig{Address: ""},
				ServerPort:  "8080",
				BaseURL:     "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "Missing Server Port",
			config: &Config{
				RedisConfig: &RedisConfig{Address: "localhost:6379"},
				ServerPort:  "",
				BaseURL:     "http://localhost:8080",
			},
			wantErr: true,
		},
		{
			name: "Missing Base URL",
			config: &Config{
				RedisConfig: &RedisConfig{Address: "localhost:6379"},
				ServerPort:  "8080",
				BaseURL:     "",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.config)

			if (err != nil) != tc.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestGetEnvFunctions(t *testing.T) {
	testCases := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue interface{}
		expected     interface{}
	}{
		{
			name:         "String Env Variable",
			envKey:       "TEST_STRING_ENV",
			envValue:     "test-value",
			defaultValue: "default-value",
			expected:     "test-value",
		},
		{
			name:         "Int Env Variable",
			envKey:       "TEST_INT_ENV",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "Duration Env Variable",
			envKey:       "TEST_DURATION_ENV",
			envValue:     "1h",
			defaultValue: 30 * time.Minute,
			expected:     1 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable
			resetEnv := setEnv(tc.envKey, tc.envValue)
			defer resetEnv()

			var result interface{}
			switch v := tc.defaultValue.(type) {
			case string:
				result = getEnv(tc.envKey, v)
			case int:
				result = getEnvAsInt(tc.envKey, v)
			case time.Duration:
				result = getEnvAsDuration(tc.envKey, v)
			}

			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// Bencmark for loading configuration
func BenchmarkLoadConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}
