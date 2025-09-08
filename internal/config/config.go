package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig   `mapstructure:"logger"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret         string `mapstructure:"access_secret"`
	RefreshSecret        string `mapstructure:"refresh_secret"`
	AccessTokenDuration  string `mapstructure:"access_token_duration"`
	RefreshTokenDuration string `mapstructure:"refresh_token_duration"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	OutputPath string `mapstructure:"output_path"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set environment variable prefix
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read from .env file if it exists
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist, we'll use env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)

	// Database defaults
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.database", "app.db")
	v.SetDefault("database.ssl_mode", "disable")

	// JWT defaults
	v.SetDefault("jwt.access_secret", "your-access-secret")
	v.SetDefault("jwt.refresh_secret", "your-refresh-secret")
	v.SetDefault("jwt.access_token_duration", "15m")
	v.SetDefault("jwt.refresh_token_duration", "7d")

	// Logger defaults
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.encoding", "json")
	v.SetDefault("logger.output_path", "stdout")
}

// GetDSN returns database connection string based on the database type
func (c *DatabaseConfig) GetDSN() string {
	switch strings.ToLower(c.Type) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			c.Host, c.Username, c.Password, c.Database, c.Port, c.SSLMode)
	case "sqlite":
		return c.Database
	default:
		return c.Database
	}
}