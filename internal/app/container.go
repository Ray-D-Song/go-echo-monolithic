package app

import (
	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/config"
	"github.com/ray-d-song/go-echo-monolithic/internal/handler"
	"github.com/ray-d-song/go-echo-monolithic/internal/middleware"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/database"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/jwt"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/validator"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"github.com/ray-d-song/go-echo-monolithic/internal/service"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Container configures and provides dependency injection
var Container = fx.Options(
	// Config
	fx.Provide(config.Load),

	// Logger
	fx.Provide(func(cfg *config.Config) (*logger.Logger, error) {
		return logger.NewLogger(&cfg.Logger)
	}),

	// Database
	fx.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		conn, err := database.NewConnection(&cfg.Database)
		if err != nil {
			return nil, err
		}
		return conn.DB, nil
	}),

	// JWT Manager
	fx.Provide(func(cfg *config.Config) (*jwt.Manager, error) {
		return jwt.NewManager(&cfg.JWT)
	}),

	// Validator
	fx.Provide(validator.NewValidator),

	// Repositories
	fx.Provide(func(db *gorm.DB) *repository.UserRepository {
		return repository.NewUserRepository(db)
	}),
	fx.Provide(func(db *gorm.DB) *repository.AuthRepository {
		return repository.NewAuthRepository(db)
	}),
	fx.Provide(func(db *gorm.DB) *repository.Migrator {
		return repository.NewMigrator(db)
	}),

	// Services
	fx.Provide(func(userRepo *repository.UserRepository) *service.UserService {
		return service.NewUserService(userRepo)
	}),
	fx.Provide(func(
		userRepo *repository.UserRepository,
		authRepo *repository.AuthRepository,
		jwtManager *jwt.Manager,
		validator *validator.Validator,
		userService *service.UserService,
	) *service.AuthService {
		return service.NewAuthService(userRepo, authRepo, jwtManager, validator, userService)
	}),
	fx.Provide(func(logger *logger.Logger) *service.WebSocketService {
		return service.NewWebSocketService(logger)
	}),

	// Handlers
	fx.Provide(func(authService *service.AuthService) *handler.AuthHandler {
		return handler.NewAuthHandler(authService)
	}),
	fx.Provide(func(userService *service.UserService) *handler.UserHandler {
		return handler.NewUserHandler(userService)
	}),
	fx.Provide(func(wsService *service.WebSocketService, logger *logger.Logger) *handler.WebSocketHandler {
		return handler.NewWebSocketHandler(wsService, logger)
	}),

	// Middleware
	fx.Provide(
		fx.Annotate(
			func(jwtManager *jwt.Manager) echo.MiddlewareFunc {
				return middleware.JWTAuth(jwtManager)
			},
			fx.ResultTags(`name:"JWTAuthMiddleware"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			func(logger *logger.Logger) echo.MiddlewareFunc {
				return middleware.Logger(logger)
			},
			fx.ResultTags(`name:"LoggerMiddleware"`),
		),
	),

	// Server
	fx.Provide(NewServer),
)
