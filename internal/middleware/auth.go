package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/jwt"
	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/response"
)

// JWTAuth returns JWT authentication middleware
func JWTAuth(jwtManager *jwt.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized(c, "Missing authorization header")
			}

			// Check if header has Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return response.Unauthorized(c, "Invalid authorization header format")
			}

			tokenString := parts[1]
			if tokenString == "" {
				return response.Unauthorized(c, "Missing token")
			}

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				return response.Unauthorized(c, "Invalid or expired token")
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// OptionalJWTAuth returns optional JWT authentication middleware
// This middleware doesn't return an error if no token is provided
func OptionalJWTAuth(jwtManager *jwt.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			// Check if header has Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return next(c)
			}

			tokenString := parts[1]
			if tokenString == "" {
				return next(c)
			}

			// Validate token
			claims, err := jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				return next(c)
			}

			// Set user information in context if token is valid
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) (uint, bool) {
	userID, ok := c.Get("user_id").(uint)
	return userID, ok
}

// GetUsername extracts username from context
func GetUsername(c echo.Context) (string, bool) {
	username, ok := c.Get("username").(string)
	return username, ok
}

// GetEmail extracts email from context
func GetEmail(c echo.Context) (string, bool) {
	email, ok := c.Get("email").(string)
	return email, ok
}