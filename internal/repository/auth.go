package repository

import (
	"errors"
	"time"

	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"github.com/ray-d-song/go-echo-monolithic/internal/types"
	"gorm.io/gorm"
)

// AuthRepository handles authentication-related data operations
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateRefreshToken creates a new refresh token
func (r *AuthRepository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetRefreshToken retrieves a refresh token by token string
func (r *AuthRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	if err := r.db.Preload("User").Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrInvalidToken
		}
		return nil, err
	}
	return &refreshToken, nil
}

// GetRefreshTokenByUserID retrieves active refresh tokens for a user
func (r *AuthRepository) GetRefreshTokenByUserID(userID uint) ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := r.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?", 
		userID, false, time.Now()).Find(&tokens).Error
	return tokens, err
}

// RevokeRefreshToken revokes a refresh token
func (r *AuthRepository) RevokeRefreshToken(token string) error {
	result := r.db.Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true)
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return types.ErrInvalidToken
	}
	
	return nil
}

// RevokeAllUserRefreshTokens revokes all refresh tokens for a user
func (r *AuthRepository) RevokeAllUserRefreshTokens(userID uint) error {
	return r.db.Model(&model.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

// CleanupExpiredTokens removes expired refresh tokens
func (r *AuthRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&model.RefreshToken{}).Error
}

// IsRefreshTokenValid checks if refresh token is valid (not revoked and not expired)
func (r *AuthRepository) IsRefreshTokenValid(token string) (bool, error) {
	var count int64
	err := r.db.Model(&model.RefreshToken{}).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", 
			token, false, time.Now()).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// UpdateRefreshToken updates an existing refresh token
func (r *AuthRepository) UpdateRefreshToken(token *model.RefreshToken) error {
	return r.db.Save(token).Error
}