// Package main provides the entry point for the go-echo-monolithic service
//
//	@title			Go Echo Monolithic API
//	@version		1.0
//	@description	A monolithic Go service built with Echo framework, featuring JWT authentication, WebSocket support, and comprehensive middleware
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ray-d-song/go-echo-monolithic/internal/app"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/ray-d-song/go-echo-monolithic/docs"
)

func main() {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	// Start application
	app := fx.New(
		app.Container,
		fx.Invoke(startServer),
	)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-ctx.Done()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}
}

// startServer starts the HTTP server
func startServer(
	server *app.Server,
	migrator *repository.Migrator,
	logger *logger.Logger,
	lifecycle fx.Lifecycle,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(ctx, migrator); err != nil {
					logger.Error("Server failed to start", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Server stopped")
			return nil
		},
	})
}
