package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/jwt"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/validator"
	"github.com/ray-d-song/go-echo-monolithic/internal/repository"
	"github.com/ray-d-song/go-echo-monolithic/internal/types"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repository.UserRepository
	authRepo    *repository.AuthRepository
	jwtManager  *jwt.Manager
	validator   *validator.Validator
	userService *UserService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repository.UserRepository,
	authRepo *repository.AuthRepository,
	jwtManager *jwt.Manager,
	validator *validator.Validator,
	userService *UserService,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		authRepo:    authRepo,
		jwtManager:  jwtManager,
		validator:   validator,
		userService: userService,
	}
}

// Register registers a new user
func (s *AuthService) Register(req *types.RegisterRequest) (*types.AuthResponse, error) {
	// Validate input
	if err := s.validator.ValidateUsername(req.Username); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := s.validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateName(req.FirstName, "first name"); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateName(req.LastName, "last name"); err != nil {
		return nil, err
	}

	// Check if user already exists
	if exists, err := s.userRepo.ExistsByUsername(req.Username); err != nil {
		return nil, err
	} else if exists {
		return nil, types.ErrUserAlreadyExists
	}

	if exists, err := s.userRepo.ExistsByEmail(req.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, types.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenDuration()),
		IsRevoked: false,
	}

	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return &types.AuthResponse{
		User:         s.userService.ToResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *types.LoginRequest) (*types.AuthResponse, error) {
	// Validate input
	if req.Username == "" {
		return nil, types.ErrValidationFailed
	}
	if req.Password == "" {
		return nil, types.ErrValidationFailed
	}

	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if err == types.ErrUserNotFound {
			return nil, types.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, types.ErrForbidden
	}

	// Verify password
	if !s.verifyPassword(user.PasswordHash, req.Password) {
		return nil, types.ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenDuration()),
		IsRevoked: false,
	}

	if err := s.authRepo.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return &types.AuthResponse{
		User:         s.userService.ToResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthService) RefreshToken(req *types.RefreshTokenRequest) (*types.TokenResponse, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, types.ErrInvalidToken
	}

	// Check if token exists and is not revoked
	storedToken, err := s.authRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if storedToken.IsRevoked {
		return nil, types.ErrTokenRevoked
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, types.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, types.ErrForbidden
	}

	// Generate new tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	if err := s.authRepo.RevokeRefreshToken(req.RefreshToken); err != nil {
		return nil, err
	}

	// Store new refresh token
	newRefreshToken := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenDuration()),
		IsRevoked: false,
	}

	if err := s.authRepo.CreateRefreshToken(newRefreshToken); err != nil {
		return nil, err
	}

	return &types.TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// Logout revokes refresh token
func (s *AuthService) Logout(refreshToken string) error {
	return s.authRepo.RevokeRefreshToken(refreshToken)
}

// LogoutAllDevices revokes all refresh tokens for a user
func (s *AuthService) LogoutAllDevices(userID uint) error {
	return s.authRepo.RevokeAllUserRefreshTokens(userID)
}

// generateRandomToken generates a random token string
func (s *AuthService) generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashPassword hashes a password using bcrypt
func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword verifies a password against its hash
func (s *AuthService) verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}