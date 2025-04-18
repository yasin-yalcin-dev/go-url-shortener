# Configuration Management

## 1. Configuration Overview
- Environment-based configuration
- Support for .env files and environment variables
- Flexible and extensible configuration system

## 2. Configuration Sources

### 2.1 Environment Variables
- Highest priority configuration method
- Overrides .env file settings
- Allows runtime configuration changes

### 2.2 .env File
- Default configuration method
- Easy management across different environments
- Supports local development

## 3. Configuration Categories

### 3.1 Redis Configuration
- `REDIS_ADDR`: Server address (default: localhost:6379)
- `REDIS_PASSWORD`: Authentication password
- `REDIS_DB`: Database number
- `REDIS_POOL_SIZE`: Connection pool size
- `REDIS_DIAL_TIMEOUT`: Connection timeout
- `REDIS_READ_TIMEOUT`: Read operation timeout
- `REDIS_WRITE_TIMEOUT`: Write operation timeout

### 3.2 Server Configuration
- `SERVER_PORT`: HTTP server listening port
- `BASE_URL`: Base URL for shortened links
- `LOG_LEVEL`: Logging verbosity level

### 3.3 URL Shortener Configuration
- `DEFAULT_URL_TTL`: Default URL expiration time
- `SHORT_ID_LENGTH`: Generated short ID length

## 4. Configuration Loading Process

### 4.1 Steps
1. Load .env file
2. Override with environment variables
3. Apply default values
4. Validate configuration

### 4.2 Validation Rules
- Required fields check
- Type conversion
- Value range validation

## 5. Environment-Specific Configurations

### 5.1 Development
- Verbose logging
- Local Redis instance
- Relaxed validation

### 5.2 Production
- Minimal logging
- Secure Redis configuration
- Strict validation

## 6. Best Practices

- Never commit sensitive information
- Use environment-specific .env files
- Use strong, unique passwords
- Regularly rotate credentials
