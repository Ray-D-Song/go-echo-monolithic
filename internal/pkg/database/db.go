package database

import (
	"fmt"
	"strings"

	"github.com/ray-d-song/go-echo-monolithic/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connection wraps the database connection
type Connection struct {
	DB *gorm.DB
}

// NewConnection creates a new database connection based on the configuration
func NewConnection(cfg *config.DatabaseConfig) (*Connection, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(cfg.Type) {
	case "mysql":
		dialector = mysql.Open(cfg.GetDSN())
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.GetDSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.GetDSN())
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}