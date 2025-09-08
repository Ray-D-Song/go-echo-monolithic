# Go Echo Monolithic

A modern monolithic web application built with Go, Echo framework, and clean architecture principles.

## Tech Stack

- **Web Framework**: [Echo](https://github.com/labstack/echo/v4) - High performance, extensible, minimalist Go web framework
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- **Configuration**: [Viper](https://github.com/spf13/viper) - Go configuration with fangs
- **WebSocket**: [Gorilla WebSocket](https://github.com/gorilla/websocket) - A WebSocket implementation for Go
- **Dependency Injection**: [Fx](https://go.uber.org/fx) - A dependency injection based application framework for Go
- **Logging**: [Zap](https://go.uber.org/zap) - Blazing fast, structured, leveled logging in Go
- **Authentication**: JWT with dual-token (access + refresh) system and silent refresh
- **ORM**: [GORM](https://gorm.io/gorm) - The fantastic ORM library for Golang
- **Database**: SQLite by default, configurable to MySQL/PostgreSQL via environment variables

## Features

- JWT authentication with dual-token system
- WebSocket support
- Request ID middleware
- Structured logging
- CORS middleware
- Rate limiting
- Graceful shutdown
- Database migrations
- Clean architecture with dependency injection

## Project Structure

```
go-echo-monolithic/
├── cmd/                        # Application entry points
│   ├── server/                 # HTTP server startup
│   │   └── main.go
│   └── cli/                    # CLI commands
│       └── main.go
├── internal/                   # Private application code
│   ├── app/                    # Application layer
│   │   ├── container.go        # DI container configuration
│   │   └── server.go           # HTTP server configuration
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── handler/                # HTTP handlers
│   │   ├── auth.go             # Authentication endpoints
│   │   ├── user.go             # User endpoints
│   │   └── websocket.go        # WebSocket endpoints
│   ├── middleware/             # Middleware
│   │   ├── auth.go             # JWT authentication middleware
│   │   ├── cors.go             # CORS middleware
│   │   ├── logger.go           # Logging middleware
│   │   └── requestid.go        # Request ID middleware
│   ├── service/                # Business logic layer
│   │   ├── auth.go             # Authentication service
│   │   ├── user.go             # User service
│   │   └── websocket.go        # WebSocket service
│   ├── repository/             # Data access layer
│   │   ├── auth.go             # Authentication data operations
│   │   ├── user.go             # User data operations
│   │   └── migration.go        # Database migrations
│   ├── model/                  # Data models
│   │   ├── user.go             # User model
│   │   ├── token.go            # Token model
│   │   └── common.go           # Common models
│   ├── pkg/                    # Internal utilities
│   │   ├── database/           # Database connection
│   │   │   └── db.go
│   │   ├── jwt/                # JWT utilities
│   │   │   └── jwt.go
│   │   ├── logger/             # Logger configuration
│   │   │   └── logger.go
│   │   ├── validator/          # Validator
│   │   │   └── validator.go
│   │   └── response/           # Unified response
│   │       └── response.go
│   └── types/                  # Type definitions
│       ├── request.go          # Request types
│       ├── response.go         # Response types
│       └── errors.go           # Error types
├── pkg/                        # Public libraries (can be imported by external projects)
├── api/                        # API definition files
│   └── openapi/                # OpenAPI specifications
├── web/                        # Frontend static files (if needed)
├── migrations/                 # Database migration files
├── scripts/                    # Script files
│   ├── build.sh
│   └── migrate.sh
├── configs/                    # Configuration files
│   ├── .env.example
│   └── config.yaml
├── docs/                       # Documentation
├── test/                       # Test files
│   ├── integration/            # Integration tests
│   └── fixtures/               # Test fixtures
├── .env                        # Environment variables (not committed to git)
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Make (optional, for build scripts)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-echo-monolithic
   ```

2. Copy environment configuration:
   ```bash
   cp configs/.env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run database migrations:
   ```bash
   go run cmd/cli/main.go migrate up
   ```

5. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080` by default.

## Configuration

The application uses environment variables for configuration. See `.env.example` for available options.

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout

### Users
- `GET /api/users/profile` - Get current user profile
- `PUT /api/users/profile` - Update current user profile
- `GET /api/users/:id` - Get user by ID
- `GET /api/users/username/:username` - Get user by username
- `GET /api/users` - List users with pagination
- `DELETE /api/users/:id` - Delete user

### WebSocket
- `WS /api/ws/connect` - WebSocket connection

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test ./test/integration/...
```

### Building

```bash
# Build the application
go build -o bin/server ./cmd/server

# Using make
make build
```

## License

This project is licensed under the MIT License.