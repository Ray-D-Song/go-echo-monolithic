package repository

import (
	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder handles database seeding operations
type Seeder struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB, logger *logger.Logger) *Seeder {
	return &Seeder{
		db:     db,
		logger: logger,
	}
}

// SeedAll runs all seeding operations
func (s *Seeder) SeedAll() error {
	if err := s.SeedUsers(); err != nil {
		return err
	}

	if err := s.SeedKVData(); err != nil {
		return err
	}

	return nil
}

// SeedUsers creates sample users
func (s *Seeder) SeedUsers() error {
	users := []*model.User{
		{
			Username:  "admin",
			Email:     "admin@example.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			IsActive:  true,
		},
		{
			Username:  "john_doe",
			Email:     "john.doe@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Role:      "user",
			IsActive:  true,
		},
		{
			Username:  "jane_smith",
			Email:     "jane.smith@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
			Role:      "user",
			IsActive:  true,
		},
		{
			Username:  "demo_user",
			Email:     "demo@example.com",
			FirstName: "Demo",
			LastName:  "User",
			Role:      "user",
			IsActive:  false,
		},
	}

	for _, user := range users {
		// Check if user already exists
		var existingUser model.User
		if err := s.db.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
			s.logger.Info("User already exists, skipping: " + user.Username)
			continue
		}

		// Hash password (default password: "password123")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)

		if err := s.db.Create(user).Error; err != nil {
			return err
		}
		s.logger.Info("Created user: " + user.Username)
	}

	return nil
}

// SeedKVData creates sample key-value data
func (s *Seeder) SeedKVData() error {
	kvPairs := []*model.KV{
		{
			Key:   "app_name",
			Value: "Go Echo Monolithic",
		},
		{
			Key:   "app_version",
			Value: "1.0.0",
		},
		{
			Key:   "maintenance_mode",
			Value: "false",
		},
		{
			Key:   "max_users_per_page",
			Value: "20",
		},
		{
			Key:   "default_user_role",
			Value: "user",
		},
		{
			Key:   "welcome_message",
			Value: "Welcome to our application!",
		},
	}

	for _, kv := range kvPairs {
		// Check if key already exists
		var existingKV model.KV
		if err := s.db.Where("key = ?", kv.Key).First(&existingKV).Error; err == nil {
			s.logger.Info("KV pair already exists, skipping: " + kv.Key)
			continue
		}

		if err := s.db.Create(kv).Error; err != nil {
			return err
		}
		s.logger.Info("Created KV pair: " + kv.Key + " = " + kv.Value)
	}

	return nil
}

// ClearData removes all seeded data (useful for testing)
func (s *Seeder) ClearData() error {
	// Delete in reverse order due to foreign key constraints
	if err := s.db.Unscoped().Delete(&model.RefreshToken{}, "1 = 1").Error; err != nil {
		return err
	}

	if err := s.db.Unscoped().Delete(&model.User{}, "1 = 1").Error; err != nil {
		return err
	}

	if err := s.db.Unscoped().Delete(&model.KV{}, "1 = 1").Error; err != nil {
		return err
	}

	s.logger.Info("Cleared all seeded data")
	return nil
}