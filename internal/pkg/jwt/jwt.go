package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ray-d-song/go-echo-monolithic/internal/config"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Manager handles JWT token operations
type Manager struct {
	accessSecret         []byte
	refreshSecret        []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewManager creates a new JWT manager
func NewManager(cfg *config.JWTConfig) (*Manager, error) {
	accessDuration, err := time.ParseDuration(cfg.AccessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token duration: %w", err)
	}

	refreshDuration, err := time.ParseDuration(cfg.RefreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token duration: %w", err)
	}

	return &Manager{
		accessSecret:         []byte(cfg.AccessSecret),
		refreshSecret:        []byte(cfg.RefreshSecret),
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

// GenerateTokenPair generates access and refresh token pair
func (m *Manager) GenerateTokenPair(userID uint, username, email string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := m.generateToken(userID, username, email, m.accessSecret, m.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := m.generateToken(userID, username, email, m.refreshSecret, m.refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateAccessToken validates access token and returns claims
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.accessSecret)
}

// ValidateRefreshToken validates refresh token and returns claims
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.refreshSecret)
}

// generateToken generates a JWT token with given parameters
func (m *Manager) generateToken(userID uint, username, email string, secret []byte, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-echo-monolithic",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// validateToken validates a JWT token and returns claims
func (m *Manager) validateToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetRefreshTokenDuration returns refresh token duration
func (m *Manager) GetRefreshTokenDuration() time.Duration {
	return m.refreshTokenDuration
}