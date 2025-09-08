# Claude Development Notes

## Code Standards

- **All comments must be written in English**
- Follow Go best practices and idiomatic code style
- Use meaningful variable and function names
- Maintain consistent formatting and indentation
- Include proper error handling

## Project Architecture

This is a monolithic Go service built with:
- Echo web framework
- Fx dependency injection
- GORM ORM
- JWT dual-token authentication
- Zap logger
- Viper configuration management
- WebSocket support
- CLI commands with Cobra

## Development Guidelines

- Use dependency injection pattern throughout the application
- Implement proper separation of concerns (handler -> service -> repository)
- Follow the established directory structure
- Include appropriate middleware for logging, request ID, and authentication
- Implement proper database migrations
- Use environment variables for configuration
- Write comprehensive error messages
- Include validation for all input data

## Testing

- Write unit tests for business logic
- Include integration tests for API endpoints
- Use test fixtures for database testing
- Mock external dependencies

## Security

- Implement secure JWT token handling
- Use proper password hashing
- Validate all input data
- Implement rate limiting
- Use CORS middleware appropriately