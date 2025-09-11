package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/ray-d-song/go-echo-monolithic/internal/config"
	"github.com/ray-d-song/go-echo-monolithic/internal/handler"
	"github.com/ray-d-song/go-echo-monolithic/internal/middleware"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"github.com/ray-d-song/go-echo-monolithic/internal/static"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/ray-d-song/go-echo-monolithic/docs"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
	logger *logger.Logger
}

// ServerParams holds server dependencies
type ServerParams struct {
	fx.In

	Config   *config.Config
	Logger   *logger.Logger
	Migrator *repository.Migrator

	// Handlers
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	WebSocketHandler *handler.WebSocketHandler
	ConfigHandler    *handler.ConfigHandler

	// Middleware
	AuthMiddleware   echo.MiddlewareFunc `name:"JWTAuthMiddleware"`
	LoggerMiddleware echo.MiddlewareFunc `name:"LoggerMiddleware"`
}

// NewServer creates a new HTTP server
func NewServer(params ServerParams) *Server {
	e := echo.New()

	// Hide Echo banner
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:   e,
		config: params.Config,
		logger: params.Logger,
	}

	// Configure middleware
	server.setupMiddleware(params)

	// Setup routes
	server.setupRoutes(params)

	return server
}

// setupMiddleware configures Echo middleware
func (s *Server) setupMiddleware(params ServerParams) {
	// Request ID middleware
	s.echo.Use(echoMiddleware.RequestID())

	// Custom logger middleware
	s.echo.Use(params.LoggerMiddleware)

	// Recovery middleware
	s.echo.Use(echoMiddleware.Recover())

	// CORS middleware
	s.echo.Use(middleware.CORS())

	// Security headers middleware
	s.echo.Use(echoMiddleware.Secure())

	// Request timeout middleware
	s.echo.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Rate limiting middleware
	s.echo.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
}

// setupRoutes configures application routes
func (s *Server) setupRoutes(params ServerParams) {
	// Swagger documentation
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check endpoint
	s.echo.GET("/health", s.healthCheck)

	// API version info
	s.echo.GET("/api/version", s.versionInfo)

	// Register handler routes
	params.AuthHandler.RegisterRoutes(s.echo)
	params.UserHandler.RegisterRoutes(s.echo, params.AuthMiddleware)
	params.WebSocketHandler.RegisterRoutes(s.echo, params.AuthMiddleware)
	params.ConfigHandler.RegisterRoutes(s.echo, params.AuthMiddleware)

	// Embedded static file serving for SPA
	s.echo.Use(echoMiddleware.StaticWithConfig(echoMiddleware.StaticConfig{
		Root:   "/",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
		Skipper: func(c echo.Context) bool {
			// Skip static file handling for API routes
			path := c.Request().URL.Path
			return (len(path) >= 4 && path[:4] == "/api") ||
				path == "/health" ||
				(len(path) >= 8 && path[:8] == "/swagger")
		},
		Filesystem: http.FS(static.GetWebFS()),
	}))
}

// healthCheck handles health check requests
// @Summary		Health check
// @Description	Check the health status of the service
// @Tags			system
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"Service is healthy"
// @Router			/health [get]
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "go-echo-monolithic",
	})
}

// versionInfo handles version info requests
// @Summary		Version info
// @Description	Get service version information
// @Tags			system
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"Version information"
// @Router			/api/version [get]
func (s *Server) versionInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version":   "1.0.0",
		"service":   "go-echo-monolithic",
		"timestamp": time.Now().UTC(),
	})
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context, migrator *repository.Migrator) error {
	// Run database migrations
	if err := migrator.AutoMigrate(); err != nil {
		s.logger.Fatal("Failed to run migrations", zap.Error(err))
		return err
	}

	// Create indexes
	if err := migrator.CreateIndexes(); err != nil {
		s.logger.Warn("Failed to create indexes", zap.Error(err))
	}

	// Start server
	addr := fmt.Sprintf(":%d", s.config.Server.Port)
	s.logger.Info("Starting HTTP server",
		zap.String("addr", addr),
	)

	// Start server in goroutine
	go func() {
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("Shutting down server...")
	return s.echo.Shutdown(shutdownCtx)
}

// GetEcho returns the Echo instance for testing
func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}
