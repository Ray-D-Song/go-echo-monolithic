package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
	"github.com/ray-d-song/go-echo-monolithic/internal/service"
	"github.com/ray-d-song/go-echo-monolithic/internal/types"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req types.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request data")
	}

	result, err := h.authService.Register(&req)
	if err != nil {
		switch err {
		case types.ErrUserAlreadyExists:
			return response.Conflict(c, "User already exists")
		case types.ErrValidationFailed:
			return response.BadRequest(c, "Validation failed")
		default:
			return response.InternalServerError(c, "Registration failed")
		}
	}

	return response.Success(c, result, "User registered successfully")
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req types.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request data")
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case types.ErrInvalidCredentials:
			return response.Unauthorized(c, "Invalid credentials")
		case types.ErrForbidden:
			return response.Forbidden(c, "Account is disabled")
		case types.ErrValidationFailed:
			return response.BadRequest(c, "Validation failed")
		default:
			return response.InternalServerError(c, "Login failed")
		}
	}

	return response.Success(c, result, "Login successful")
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req types.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request data")
	}

	result, err := h.authService.RefreshToken(&req)
	if err != nil {
		switch err {
		case types.ErrInvalidToken:
			return response.Unauthorized(c, "Invalid refresh token")
		case types.ErrTokenExpired:
			return response.Unauthorized(c, "Refresh token expired")
		case types.ErrTokenRevoked:
			return response.Unauthorized(c, "Refresh token revoked")
		case types.ErrForbidden:
			return response.Forbidden(c, "Account is disabled")
		default:
			return response.InternalServerError(c, "Token refresh failed")
		}
	}

	return response.Success(c, result, "Token refreshed successfully")
}

// Logout handles user logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req types.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request data")
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		return response.InternalServerError(c, "Logout failed")
	}

	return response.Success(c, nil, "Logout successful")
}

// LogoutAllDevices handles logout from all devices
func (h *AuthHandler) LogoutAllDevices(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	if err := h.authService.LogoutAllDevices(userID); err != nil {
		return response.InternalServerError(c, "Logout from all devices failed")
	}

	return response.Success(c, nil, "Logged out from all devices successfully")
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	auth := e.Group("/api/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", h.Logout)
	auth.POST("/logout-all", h.LogoutAllDevices)
}
