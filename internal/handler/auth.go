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
// @Summary		Register a new user
// @Description	Register a new user with username, email and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		types.RegisterRequest	true	"Registration request"
// @Success		200		{object}	response.Response		"User registered successfully"
// @Failure		400		{object}	response.Response		"Bad request"
// @Failure		409		{object}	response.Response		"User already exists"
// @Failure		500		{object}	response.Response		"Internal server error"
// @Router			/auth/register [post]
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
// @Summary		User login
// @Description	Authenticate user with username/email and password
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		types.LoginRequest	true	"Login request"
// @Success		200		{object}	response.Response	"Login successful"
// @Failure		400		{object}	response.Response	"Bad request"
// @Failure		401		{object}	response.Response	"Invalid credentials"
// @Failure		500		{object}	response.Response	"Internal server error"
// @Router			/auth/login [post]
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
// @Summary		Refresh access token
// @Description	Refresh access token using refresh token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		types.RefreshTokenRequest	true	"Refresh token request"
// @Success		200		{object}	response.Response			"Token refreshed successfully"
// @Failure		400		{object}	response.Response			"Bad request"
// @Failure		401		{object}	response.Response			"Invalid or expired token"
// @Failure		403		{object}	response.Response			"Account disabled"
// @Failure		500		{object}	response.Response			"Internal server error"
// @Router			/auth/refresh [post]
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
// @Summary		User logout
// @Description	Logout user and revoke refresh token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		types.LogoutRequest	true	"Logout request"
// @Success		200		{object}	response.Response	"Logout successful"
// @Failure		400		{object}	response.Response	"Bad request"
// @Failure		500		{object}	response.Response	"Internal server error"
// @Router			/auth/logout [post]
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
// @Summary		Logout from all devices
// @Description	Logout user from all devices by revoking all refresh tokens
// @Tags			auth
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200		{object}	response.Response	"Logged out from all devices successfully"
// @Failure		401		{object}	response.Response	"Unauthorized"
// @Failure		500		{object}	response.Response	"Internal server error"
// @Router			/auth/logout-all [post]
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
