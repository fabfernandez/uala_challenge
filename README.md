# Uala Challenge - Microblogging API

A simplified microblogging platform API similar to Twitter, built in Go with Clean Architecture.

## Features

- **Tweets**: Post short messages (max 280 characters)
- **Follow**: Follow/unfollow other users
- **Timeline**: View tweets from users you follow
- **User Management**: Basic user identification via headers

## Quick Start

### Docker (Recommended)
```bash
docker-compose up --build
```

### Local Development
```bash
go mod tidy
go run .
```

Server starts on `http://localhost:8080`

## API Endpoints

All endpoints require `X-User-ID` header for user identification.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tweets` | Create a tweet |
| GET | `/api/v1/timeline` | Get timeline of followed users' tweets |
| GET | `/api/v1/users/tweets?user_id={id}` | Get specific user's tweets |
| POST | `/api/v1/follow` | Follow a user |
| POST | `/api/v1/unfollow` | Unfollow a user |
| GET | `/api/v1/health` | Health check |

### Example Usage

**Create a tweet:**
```bash
curl -X POST http://localhost:8080/api/v1/tweets \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{"content": "Hello, world!"}'
```

**Follow a user:**
```bash
curl -X POST http://localhost:8080/api/v1/follow \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{"followee_id": "user456"}'
```

**Get timeline:**
```bash
curl -X GET http://localhost:8080/api/v1/timeline \
  -H "X-User-ID: user123"
```

## Architecture

Built with **Clean Architecture** principles:

```
Interface Layer → Application Layer → Domain Layer
       ↓                ↓
Infrastructure Layer ←──────┘
```

### Layers

- **Domain** (`internal/domain/`): Core business entities and rules
- **Application** (`internal/application/`): Use cases and business logic
- **Infrastructure** (`internal/infrastructure/`): Storage and external services
- **Interface** (`internal/interfaces/`): HTTP handlers and routing

### Key Design Decisions

- **In-Memory Storage**: Thread-safe storage for boilerplate (production would use database)
- **User Identification**: `X-User-ID` header as per requirements
- **Character Limit**: 280 characters per tweet
- **Dependency Injection**: Use cases receive dependencies through interfaces
- **Thread Safety**: All storage operations protected with mutex locks

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker

```bash
# Build and run
docker-compose up --build

# Build image
docker build -t uala-microblog-api .

# Run container
docker run -p 8080:8080 uala-microblog-api

# Stop
docker-compose down
```

## Project Structure

```
├── main.go                    # Application entry point
├── internal/
│   ├── domain/               # Core business entities
│   ├── application/usecases/ # Business logic
│   ├── infrastructure/       # Storage implementations
│   └── interfaces/http/      # HTTP handlers
├── Dockerfile
├── docker-compose.yml
└── go.mod
```
