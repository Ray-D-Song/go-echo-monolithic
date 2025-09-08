package repository

import (
	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"gorm.io/gorm"
)

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrate runs automatic migrations for all models
func (m *Migrator) AutoMigrate() error {
	return m.db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
}

// DropTables drops all tables (use with caution)
func (m *Migrator) DropTables() error {
	return m.db.Migrator().DropTable(
		&model.RefreshToken{},
		&model.User{},
	)
}

// CreateIndexes creates additional database indexes
func (m *Migrator) CreateIndexes() error {
	// Create composite indexes for better query performance
	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_valid ON refresh_tokens(user_id, is_revoked, expires_at)").Error; err != nil {
		return err
	}

	if err := m.db.Exec("CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active)").Error; err != nil {
		return err
	}

	return nil
}