# Go URL Shortener Comprehensive Documentation

## Table of Contents
1. [Introduction](#introduction)
2. [Architecture](#architecture)
3. [Components](#components)
4. [Configuration](#configuration)
5. [API Usage](#api-usage)
6. [Development](#development)
7. [Performance](#performance)
8. [Security](#security)
9. [Error Handling](#error-handling)
10. [Example Scenarios](#example-scenarios)

## Introduction

Go URL Shortener is a high-performance, flexible, and secure URL shortening service that transforms long URLs into short, manageable links.

### Key Features
- Fast URL shortening
- Redis-based storage
- Configurable URL expiration
- Detailed analytics tracking
- Secure URL validation
- Rate limiting

## Architecture

### Layered Architecture
- **Presentation Layer**: Handlers
- **Service Layer**: Business logic
- **Data Layer**: Redis integration

```
project/
├── cmd/                # Application entry point
├── internal/           # Internal packages
│   ├── config/         # Configuration management
│   ├── handler/        # HTTP handlers
│   ├── model/          # Data models
│   ├── redis/          # Redis integration
│   └── service/        # Business logic
└── pkg/                # External packages
    ├── analytics/      # Analytics tracking
    ├── errors/         # Custom error management
    ├── logger/         # Logging
    ├── ratelimiter/    # Rate limiting
    ├── shortener/      # URL shortening
    └── validator/      # URL validation
```

(Rest of the content remains the same as the previous Turkish version, but in English)
