package model

import (
	"time"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	BaseModel
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}