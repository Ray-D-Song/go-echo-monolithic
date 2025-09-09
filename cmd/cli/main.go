package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ray-d-song/go-echo-monolithic/internal/app"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI tool for go-echo-monolithic service",
	Long:  "Command line interface for managing the go-echo-monolithic service",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  "Run database migrations to set up or update the database schema",
	Run: func(cmd *cobra.Command, args []string) {
		runWithDI(func(migrator *repository.Migrator, logger *logger.Logger) {
			logger.Info("Running database migrations...")
			
			if err := migrator.AutoMigrate(); err != nil {
				logger.Fatal("Migration failed", zap.Error(err))
				return
			}
			
			if err := migrator.CreateIndexes(); err != nil {
				logger.Warn("Failed to create indexes", zap.Error(err))
			}
			
			logger.Info("Migration completed successfully")
		})
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback database migrations",
	Long:  "Drop all database tables (use with caution)",
	Run: func(cmd *cobra.Command, args []string) {
		runWithDI(func(migrator *repository.Migrator, logger *logger.Logger) {
			logger.Warn("Rolling back database migrations (dropping all tables)...")
			
			if err := migrator.DropTables(); err != nil {
				logger.Fatal("Rollback failed", zap.Error(err))
				return
			}
			
			logger.Info("Rollback completed successfully")
		})
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup expired tokens",
	Long:  "Remove expired refresh tokens from the database",
	Run: func(cmd *cobra.Command, args []string) {
		runWithDI(func(authRepo *repository.AuthRepository, logger *logger.Logger) {
			logger.Info("Cleaning up expired tokens...")
			
			if err := authRepo.CleanupExpiredTokens(); err != nil {
				logger.Error("Cleanup failed", zap.Error(err))
				return
			}
			
			logger.Info("Token cleanup completed successfully")
		})
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the current version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-echo-monolithic CLI v1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runWithDI runs a function with dependency injection
func runWithDI(fn interface{}) {
	app := fx.New(
		app.Container,
		fx.Invoke(fn),
		fx.NopLogger, // Suppress fx logs for CLI
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}