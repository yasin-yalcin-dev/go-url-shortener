# Go URL Shortener 🔗✨

## Project Description

Go URL Shortener is a high-performance, flexible URL shortening service developed using Go and Redis. It transforms long URLs into short, manageable links.

## Features

- 🚀 Fast and lightweight URL shortening
- 🔒 Secure URL validation
- ⏱️ Configurable URL expiration
- 📊 Detailed analytics tracking
- 🛡️ Rate limiting
- 🔍 Custom short ID generation

## Requirements

- Go 1.24.2+
- Redis

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yasin-yalcin-dev/go-url-shortener.git
cd go-url-shortener
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure the .env file:
```bash
cp .env.example .env
```

4. Edit environment variables:
```
REDIS_ADDR=localhost:6379
SERVER_PORT=8080
BASE_URL=http://localhost:8080
LOG_LEVEL=info
DEFAULT_URL_TTL=24h
```

## Running the Application

```bash
go run cmd/main.go
```

## API Usage

### Shorten URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"original":"https://example.com"}'
```

### Shorten URL with Custom Expiration

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"original":"https://example.com", "ttl":3600000000000}'
```

### Get Analytics

```bash
curl http://localhost:8080/abc123/analytics
```

## Configuration

All configurations can be made via the .env file or environment variables.

### Redis Settings
- `REDIS_ADDR`: Redis server address
- `REDIS_PASSWORD`: Redis password
- `REDIS_DB`: Database to use
- `REDIS_POOL_SIZE`: Connection pool size

### Server Settings
- `SERVER_PORT`: Server port
- `BASE_URL`: Base URL
- `LOG_LEVEL`: Logging level
- `DEFAULT_URL_TTL`: Default URL expiration

## Development

- Run tests: `go test ./...`
- Format code: `go fmt ./...`

## Architecture

- Layered architecture
- Dependency injection
- Modular design

## Contributing

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Submit a pull request

## License

MIT License

